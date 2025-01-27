package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
	"github.com/roadmap-thesis/backend/pkg/llm"
	"github.com/rs/zerolog/log"
)

func (b *backend) GenerateRoadmap(ctx context.Context, input io.GenerateRoadmapInput) (io.GenerateRoadmapOutput, error) {
	var output io.GenerateRoadmapOutput

	generated, err := b.chatGeneratePrompt(ctx, llm.ChatPrompt{
		System: b.makeGenerateRoadmapSystemPrompt(),
		User:   b.makeGenerateRoadmapUserPrompt(input),
	})
	if err != nil {
		return io.GenerateRoadmapOutput{}, err
	}

	auth := auth.FromContext(ctx)
	roadmap := domain.NewRoadmap(auth.ID, generated.Title, generated.Description)

	for _, topic := range generated.Topics {
		newTopic := domain.NewTopic(topic.Title, topic.Description)
		roadmap.AddTopic(newTopic)
		if len(topic.Subtopics) > 0 {
			for _, subtopic := range topic.Subtopics {
				newSubtopic := domain.NewTopic(subtopic.Title, subtopic.Description)
				newTopic.AddSubtopic(newSubtopic)
			}
		}
	}

	createdRoadmap, err := b.repository.Roadmap.Save(ctx, roadmap)
	if err != nil {
		return io.GenerateRoadmapOutput{}, err
	}

	output.Slug = createdRoadmap.Slug

	return output, nil
}

type chatGeneratePromptPromptResult struct {
	Title       string                                `json:"title"`
	Description string                                `json:"description"`
	Topics      []chatGeneratePromptPromptResultTopic `json:"topics"`
}

type chatGeneratePromptPromptResultTopic struct {
	Title       string                                   `json:"title"`
	Description string                                   `json:"description"`
	Subtopics   []chatGeneratePromptPromptResultSubtopic `json:"subtopics"`
}

type chatGeneratePromptPromptResultSubtopic struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (b *backend) chatGeneratePrompt(ctx context.Context, prompt llm.ChatPrompt) (chatGeneratePromptPromptResult, error) {
	fmt.Println(prompt.User)
	content, err := b.llm.Chat(ctx, prompt)
	if err != nil {
		return chatGeneratePromptPromptResult{}, err
	}

	var result chatGeneratePromptPromptResult

	log.Debug().Msg(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return chatGeneratePromptPromptResult{}, err
	}

	return result, nil
}

func (b *backend) makeGenerateRoadmapUserPrompt(input io.GenerateRoadmapInput) string {
	var sb strings.Builder
	sb.WriteString(`I will give you a topic and you need to generate a learning roadmap for it. Just reply to the question without adding any other information about the prompt and use simple language.
`)

	sb.WriteString("Generate a structured learning roadmap for the topic: ")
	sb.WriteString(input.Topic)

	sb.WriteString("\nHere are my personalization options:\n")

	sb.WriteString("- Daily Time Availability: ")
	sb.WriteString(input.PersonalizationOptions.DailyTimeAvailability)
	sb.WriteString("\n")

	sb.WriteString("- Total Duration: ")
	sb.WriteString(input.PersonalizationOptions.TotalDuration)
	sb.WriteString("\n")

	sb.WriteString("- Skill Level: ")
	sb.WriteString(string(input.PersonalizationOptions.SkillLevel))
	sb.WriteString("\n")

	if input.PersonalizationOptions.AdditionalInfo != "" {
		sb.WriteString("- Additional Information:\n \"\"\"\n ")
		sb.WriteString(input.PersonalizationOptions.AdditionalInfo)
		sb.WriteString("\n \"\"\"\n")
	}

	return sb.String()
}

func (b *backend) makeGenerateRoadmapSystemPrompt() string {
	var sb strings.Builder

	promptUserPersonalizationOptions := []string{
		"Daily Time Availability: How much time the user can dedicate daily (e.g., 15 minutes, 30 minutes, 1 hour).",
		"Total Duration: The overall duration of the roadmap (e.g., 1 week, 3 months).",
		"Skill Level: The user's experience level (e.g., Beginner, Intermediate, Advanced).",
		"Additional Info: Any other user-provided goals or preferences. This is Optional for the user.",
	}

	promptSystemGuidelines := []string{
		"Go into detail about the main topic to provide a comprehensive overview of the subject.",
		"Each topic should have a title and a brief description to explain the focus of that section.",
		"Subtopics should be related to the main topic and provide more detailed information on specific aspects of the subject.",
		"Each description should be clear and informative. It should be long enough to explain the topic but concise enough to maintain the user's interest.",
		"Ensure that a topic is broken down into manageable subtopics to help users understand the subject better whenever possible.",
		"A topic can also not have any subtopics if it is a standalone subject.",
		"Ensure the roadmap is concise, user-friendly, and structured for easy navigation.",
		"Use only English language for the roadmap.",
		"Use plain and neutral language suitable for learners of all backgrounds.",
	}

	exampleFormat := chatGeneratePromptPromptResult{
		Title:       "Example Topic",
		Description: "An extensive overview of the topic to set the stage for learning.",
		Topics: []chatGeneratePromptPromptResultTopic{
			{
				Title:       "Main Topic",
				Description: "A one paragraph long explanation of the main topic.",
				Subtopics: []chatGeneratePromptPromptResultSubtopic{
					{
						Title:       "Subtopic 1",
						Description: "A one paragraph long explanation of Subtopic 1.",
					},
					{
						Title:       "Subtopic 2",
						Description: "A one paragraph long explanation of Subtopic 2.",
					},
				},
			},
		},
	}

	exampleResult := chatGeneratePromptPromptResult{
		Title:       "Front End Development",
		Description: "Step by step guide to learn  frontend development.",
		Topics: []chatGeneratePromptPromptResultTopic{
			{
				Title:       "What Is Front End Dev?",
				Description: "Front end development is the practice of producing HTML, CSS, and JavaScript for a website or web application so a user can see and interact with them directly. It involves the design of the site, the layout, the colors, the fonts, and so on.",
				Subtopics: []chatGeneratePromptPromptResultSubtopic{
					{
						Title:       "HTML",
						Description: "HTML is the standard markup language for creating web pages and web applications. It provides the basic structure of sites, which is enhanced and modified by other technologies like CSS and JavaScript.",
					},
					{
						Title:       "CSS",
						Description: "CSS is a style sheet language used for describing the presentation of a document written in HTML. It controls the layout of multiple web pages all at once.",
					},
					{
						Title:       "JavaScript",
						Description: "JavaScript is a programming language that enables you to interact with elements on a webpage. It is used for creating dynamic and interactive web pages.",
					},
					{
						Title:       "Responsive Design",
						Description: "Responsive design is an approach to web design that makes web pages render well on a variety of devices and window or screen sizes.",
					},
				},
			},
			{
				Title:       "JavaScript Frameworks and Libraries",
				Description: "JavaScript frameworks and libraries are pre-written JavaScript code that helps you build interactive web applications. They provide ready-to-use functions and components that you can use in your code.",
				Subtopics: []chatGeneratePromptPromptResultSubtopic{
					{
						Title:       "React",
						Description: "React is a JavaScript library for building user interfaces. It is maintained by Facebook and a community of individual developers and companies.",
					},
					{
						Title:       "Vue.js",
						Description: "Vue.js is a progressive JavaScript framework used to build interactive web interfaces. It is designed from the ground up to be incrementally adoptable.",
					},
					{
						Title:       "Angular",
						Description: "Angular is a platform and framework for building single-page client applications using HTML and TypeScript. It is maintained by Google.",
					},
					{
						Title:       "Svelte",
						Description: "Svelte is a radical new approach to building user interfaces. It shifts the work of rendering from the browser to the compile step, resulting in faster load times and a better user experience.",
					},
					{
						Title:       "Node.js",
						Description: "Node.js is an open-source, cross-platform, JavaScript runtime environment that executes JavaScript code outside a web browser. It is used to build scalable network applications.",
					},
				},
			},
		},
	}

	sb.WriteString(`You are an expert in creating structured learning roadmaps for a learning application. The roadmaps you generate are designed to provide users with a clear and organized path for self-learning, not as a course or detailed content provider. The roadmap will:
1. Include a title and description of the main topic to introduce the subject.
2. Break down the topic into topics and subtopics, each with a title and description to explain the focus of the section.
3. Use a maximum of 2 levels of depth for subtopics. Topics can contain subtopics, but subtopics cannot have further nested levels.
4. Be tailored based on user-provided personalization options:
`)

	for _, userPersonalizationOpt := range promptUserPersonalizationOptions {
		sb.WriteString(" - ")
		sb.WriteString(userPersonalizationOpt)
		sb.WriteString("\n")
	}

	sb.WriteString("# Guidelines:\n")
	for _, guideline := range promptSystemGuidelines {
		sb.WriteString(" - ")
		sb.WriteString(guideline)
		sb.WriteString("\n")
	}

	sb.WriteString("# Example Format:\n")

	exampleFormatJson, err := json.MarshalIndent(exampleFormat, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal example format")
		return ""
	}

	sb.Write(exampleFormatJson)

	sb.WriteString("\n")
	sb.WriteString("The roadmap must adhere to this format while reflecting the user's provided topic and personalization preferences. Do not use markdown symbols such as the triple backticks or quotes, you must only respond with the raw json itself.")
	sb.WriteString("\n")

	sb.WriteString("# A Real Case Example:\n")

	sb.WriteString("### Input:\n")
	sb.WriteString("- Topic: Front End Development\n")
	sb.WriteString("- Daily Time Availability: 1 hour/day\n")
	sb.WriteString("- Total Duration: 1 month\n")
	sb.WriteString("- Skill Level: beginner\n")

	sb.WriteString("\n### Output:\n\n")

	exampleResultJson, err := json.MarshalIndent(exampleResult, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal example format")
		return ""
	}

	sb.Write(exampleResultJson)

	return sb.String()
}

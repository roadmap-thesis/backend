package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/auth"
	"github.com/roadmap-thesis/backend/pkg/openai"
	"github.com/rs/zerolog/log"
)

func (b *backend) GenerateRoadmap(ctx context.Context, input io.GenerateRoadmapInput) (io.GenerateRoadmapOutput, error) {
	var output io.GenerateRoadmapOutput

	generated, err := b.openAiChatGenerateRoadmap(ctx, openai.ChatPrompt{
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

type openAiChatGenerateRoadmapPromptResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Topics      []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Subtopics   []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"subtopics"`
	} `json:"topics"`
}

func (b *backend) openAiChatGenerateRoadmap(ctx context.Context, prompt openai.ChatPrompt) (openAiChatGenerateRoadmapPromptResult, error) {
	fmt.Println(prompt.User)
	completionResponse, err := b.openai.Chat(ctx, prompt)
	if err != nil {
		return openAiChatGenerateRoadmapPromptResult{}, err
	}

	var result openAiChatGenerateRoadmapPromptResult
	content := completionResponse.Choices[0].Message.Content

	log.Debug().Msg(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return openAiChatGenerateRoadmapPromptResult{}, err
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

	if input.PersonalizationOptions.LearningGoal != "" {
		sb.WriteString("- Learning Goal: ")
		sb.WriteString(string(input.PersonalizationOptions.LearningGoal))
		sb.WriteString("\n")
	}

	if input.PersonalizationOptions.AdditionalInfo != "" {
		sb.WriteString("- Additional Information:\n \"\"\"\n ")
		sb.WriteString(input.PersonalizationOptions.AdditionalInfo)
		sb.WriteString("\n \"\"\"\n")
	}

	return sb.String()
}

func (b *backend) makeGenerateRoadmapSystemPrompt() string {
	return `You are an expert in creating structured learning roadmaps for a learning application. The roadmaps you generate are designed to provide users with a clear and organized path for self-learning, not as a course or detailed content provider. The roadmap will:
1. Include a title and description of the main topic to introduce the subject.
2. Break down the topic into topics and subtopics, each with a title and description to explain the focus of the section.
3. Use a maximum of 2 levels of depth for subtopics. Topics can contain subtopics, but subtopics cannot have further nested levels.
4. Be tailored based on user-provided personalization options:
   - Daily Time Availability: How much time the user can dedicate daily (e.g., 15 minutes, 30 minutes, 1 hour).
   - Total Duration: The overall duration of the roadmap (e.g., 1 week, 3 months).
   - Skill Level: The user's experience level (e.g., Beginner, Intermediate, Advanced).
   - Learning Goal: The user's learning goal for the roadmap or target (e.g., Academic, Professional, Personal). IMPORTANT: This is Optional for the user, but if provided, make it the main focus goal of the roadmap.
   - Additional Info: Any other user-provided goals or preferences. This is Optional for the user.

### Guidelines:
- Go into detail about the main topic to provide a comprehensive overview of the subject.
- Each topic should have a title and a brief description to explain the focus of that section.
- Subtopics should be related to the main topic and provide more detailed information on specific aspects of the subject.
- Each description should be clear and informative. It should be long enough to explain the topic but concise enough to maintain the user's interest.
- Ensure that a topic is broken down into manageable subtopics to help users understand the subject better whenever possible.
- A topic can also not have any subtopics if it is a standalone subject.
- Ensure the roadmap is concise, user-friendly, and structured for easy navigation.
- Use only English language for the roadmap.
- Use plain and neutral language suitable for learners of all backgrounds.

# Example Format:

"""
{
  "title": "Example Topic",
  "description": "An extensive overview of the topic to set the stage for learning.",
  "topics": [
    {
      "title": "Main Topic",
      "description": "A one paragraph long explanation of the main topic.",
      "subtopics": [
        {
          "title": "Subtopic 1",
          "description": "A one paragraph long explanation of Subtopic 1."
        },
        {
          "title": "Subtopic 2",
          "description": "A one paragraph long explanation of Subtopic 2."
        }
      ]
    }
  ]
}
"""

The roadmap must adhere to this format while reflecting the user's provided topic and personalization preferences. Do not use markdown symbols such as the triple backticks, you must only respond with the raw json itself.

# A Real Case Example:

### Input:
- Topic: Front End Development
- Daily Time Availability: 1 hour/day
- Total Duration: 1 month
- Skill Level: beginner
- Learning Goal: academic

### Output:

"""
{
  "title": "Front End Development",
  "description": "Step by step guide to learn  frontend development.",
  "topics": [
    {
      "title": "What Is Front End Dev?",
      "description": "Front end development is the practice of producing HTML, CSS, and JavaScript for a website or web application so a user can see and interact with them directly. It involves the design of the site, the layout, the colors, the fonts, and so on.",
      "subtopics": [
        {
          "title": "HTML",
		  "description": "HTML is the standard markup language for creating web pages and web applications. It provides the basic structure of sites, which is enhanced and modified by other technologies like CSS and JavaScript."
        },
        {
          "title": "CSS",
		  "description": "CSS is a style sheet language used for describing the presentation of a document written in HTML. It controls the layout of multiple web pages all at once."
        },
		{
		  "title": "JavaScript",
		  "description": "JavaScript is a programming language that enables you to interact with elements on a webpage. It is used for creating dynamic and interactive web pages."
		},
		{
		  "title": "Responsive Design",
		  "description": "Responsive design is an approach to web design that makes web pages render well on a variety of devices and window or screen sizes."
		}
	  ]
    },
	{
	  "title": "JavaScript Frameworks and Libraries",
	  "description": "JavaScript frameworks and libraries are pre-written JavaScript code that helps you build interactive web applications. They provide ready-to-use functions and components that you can use in your code.",
	  "subtopics": [
		{
		  "title": "React",
		  "description": "React is a JavaScript library for building user interfaces. It is maintained by Facebook and a community of individual developers and companies."
		},
		{
		  "title": "Vue.js",
		  "description": "Vue.js is a progressive JavaScript framework used to build interactive web interfaces. It is designed from the ground up to be incrementally adoptable."
		},
		{
		  "title": "Angular",
		  "description": "Angular is a platform and framework for building single-page client applications using HTML and TypeScript. It is maintained by Google."
		},
		{
		  "title": "Svelte",
		  "description": "Svelte is a radical new approach to building user interfaces. It shifts the work of rendering from the browser to the compile step, resulting in faster load times and a better user experience."
		},
		{
		  "title": "Node.js",
		  "description": "Node.js is an open-source, cross-platform, JavaScript runtime environment that executes JavaScript code outside a web browser. It is used to build scalable network applications."
		}
	  ]
	},
	... more topics
  ]
}
"""
`
}

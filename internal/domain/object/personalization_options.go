package object

type PersonalizationOptions struct {
	DailyTimeAvailability string       `json:"daily_time_availability" validate:"required"`
	TotalDuration         string       `json:"total_duration" validate:"required"`
	SkillLevel            SkillLevel   `json:"skill_level" validate:"required,oneof=beginner intermediate advanced"`
	LearningGoal          LearningGoal `json:"learning_goal" validate:"omitempty,oneof=academic professional personal"`
	AdditionalInfo        string       `json:"additional_info" validate:"omitempty"`
}

type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "beginner"
	SkillLevelIntermediate SkillLevel = "intermediate"
	SkillLevelAdvanced     SkillLevel = "advanced"
)

type LearningGoal string

const (
	LearningGoalAcademic     LearningGoal = "academic"
	LearningGoalProfessional LearningGoal = "professional"
	LearningGoalPersonal     LearningGoal = "personal"
)

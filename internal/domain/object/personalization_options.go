package object

type PersonalizationOptions struct {
	DailyTimeAvailability string     `json:"daily_time_availability" validate:"required"`
	TotalDuration         string     `json:"total_duration" validate:"required"`
	SkillLevel            SkillLevel `json:"skill_level" validate:"required,oneof=beginner intermediate advanced"`
	AdditionalInfo        string     `json:"additional_info" validate:"omitempty"`
}

type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "beginner"
	SkillLevelIntermediate SkillLevel = "intermediate"
	SkillLevelAdvanced     SkillLevel = "advanced"
)

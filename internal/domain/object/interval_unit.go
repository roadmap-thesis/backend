package object

type IntervalUnit string

const (
	IntervalUnitMinutes IntervalUnit = "minutes"
	IntervalUnitHours   IntervalUnit = "hours"
	IntervalUnitDays    IntervalUnit = "days"
	IntervalUnitWeeks   IntervalUnit = "weeks"
	IntervalUnitMonths  IntervalUnit = "months"
)

func (u IntervalUnit) String() string {
	return string(u)
}

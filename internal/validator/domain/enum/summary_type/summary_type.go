package summary_type

type SummaryType string

const (
	Day   SummaryType = "DAY"
	Month SummaryType = "MONTH"
	Year  SummaryType = "YEAR"
)

func (s SummaryType) Name() string {
	return string(s)
}

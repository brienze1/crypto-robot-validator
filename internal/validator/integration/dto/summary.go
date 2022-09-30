package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/summary_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
)

// Summary DynamoDB entity for crypto-robot.client repository
type Summary struct {
	Type         summary_type.SummaryType `dynamodbav:"type"`
	Day          int                      `dynamodbav:"day"`
	Month        int                      `dynamodbav:"month"`
	Year         int                      `dynamodbav:"year"`
	AmountSold   float64                  `dynamodbav:"amount_sold"`
	AmountBought float64                  `dynamodbav:"amount_bought"`
	Profit       float64                  `dynamodbav:"profit"`
}

// SummaryDto creates a dto.Summary from model.Summary
func SummaryDto(summaries []*model.Summary) []*Summary {
	var summariesDto []*Summary
	for _, summary := range summaries {
		summariesDto = append(summariesDto, &Summary{
			Type:         summary.Type,
			Day:          summary.Day,
			Month:        summary.Month,
			Year:         summary.Year,
			AmountSold:   summary.AmountSold,
			AmountBought: summary.AmountBought,
			Profit:       summary.Profit,
		})
	}

	return summariesDto
}

// ToModel creates a model.Summary from dto.Summary
func (s *Summary) ToModel() *model.Summary {
	return &model.Summary{
		Type:         s.Type,
		Day:          s.Day,
		Month:        s.Month,
		Year:         s.Year,
		AmountSold:   s.AmountSold,
		AmountBought: s.AmountBought,
		Profit:       s.Profit,
	}
}

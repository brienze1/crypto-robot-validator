package model

import "github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/summary_type"

type Summary struct {
	Type         summary_type.SummaryType
	Day          int
	Month        int
	Year         int
	AmountSold   float64
	AmountBought float64
	Profit       float64
}

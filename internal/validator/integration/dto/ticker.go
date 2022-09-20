package dto

import "strconv"

type Ticker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (t Ticker) GetPrice() (float64, error) {
	return strconv.ParseFloat(t.Price, 64)
}

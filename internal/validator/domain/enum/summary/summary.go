package summary

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/operation_type"

type Summary string

const (
	StrongBuy  Summary = "STRONG_BUY"
	Buy        Summary = "BUY"
	Neutral    Summary = "NEUTRAL"
	Sell       Summary = "SELL"
	StrongSell Summary = "STRONG_SELL"
)

var values = map[Summary]int{
	StrongBuy:  2,
	Buy:        1,
	Neutral:    0,
	Sell:       1,
	StrongSell: 2,
}

func (s Summary) Value() int {
	return values[s]
}

func (s Summary) Name() string {
	return string(s)
}

func (s Summary) OperationType() operation_type.OperationType {
	if s.Name() == Sell.Name() || s.Name() == StrongSell.Name() {
		return operation_type.Sell
	}
	if s.Name() == Buy.Name() || s.Name() == StrongBuy.Name() {
		return operation_type.Buy
	}
	return operation_type.None
}

func (s Summary) OperationTypeString() string {
	return string(s.OperationType())
}

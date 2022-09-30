package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
)

// Client DynamoDB entity for crypto-robot.client repository
type Client struct {
	Id                        string     `dynamodbav:"client_id"`
	Active                    bool       `dynamodbav:"active"`
	LockedUntil               string     `dynamodbav:"locked_until"`
	Locked                    bool       `dynamodbav:"locked"`
	CashAvailable             float64    `dynamodbav:"cash_available"`
	CashAmount                float64    `dynamodbav:"cash_amount"`
	CashReserved              float64    `dynamodbav:"cash_reserved"`
	CryptoAvailable           float64    `dynamodbav:"crypto_available"`
	CryptoAmount              float64    `dynamodbav:"crypto_amount"`
	CryptoReserved            float64    `dynamodbav:"crypto_reserved"`
	OperationStopLoss         float64    `dynamodbav:"operation_stop_loss"`
	DayStopLoss               float64    `dynamodbav:"day_stop_loss"`
	MonthStopLoss             float64    `dynamodbav:"month_stop_loss"`
	OperationAmountPercentage float64    `dynamodbav:"operation_amount_percentage"`
	BuyOn                     int        `dynamodbav:"buy_on"`
	SellOn                    int        `dynamodbav:"sell_on"`
	Symbols                   []string   `dynamodbav:"symbols"`
	Summary                   []*Summary `dynamodbav:"summary"`
}

// ClientDto creates a dto.Client from model.Client
func ClientDto(client *model.Client) *Client {
	return &Client{
		Id:                        client.Id,
		Active:                    client.Active,
		LockedUntil:               client.LockedUntil.String(),
		Locked:                    client.Locked,
		CashAvailable:             client.CashAvailable,
		CashAmount:                client.CashAmount,
		CashReserved:              client.CashReserved,
		CryptoAvailable:           client.CryptoAvailable,
		CryptoAmount:              client.CryptoAmount,
		CryptoReserved:            client.CryptoReserved,
		OperationStopLoss:         client.OperationStopLoss,
		DayStopLoss:               client.DayStopLoss,
		MonthStopLoss:             client.MonthStopLoss,
		OperationAmountPercentage: client.OperationAmountPercentage,
		BuyOn:                     client.BuyOn,
		SellOn:                    client.SellOn,
		Symbols:                   client.Symbols,
		Summary:                   SummaryDto(client.Summary),
	}
}

// ToModel creates a model.Client from dto.Client
func (client Client) ToModel() *model.Client {
	lockedUntil := time_utils.From(client.LockedUntil)

	var summaries []*model.Summary
	for _, summaryDto := range client.Summary {
		summaries = append(summaries, summaryDto.ToModel())
	}

	return &model.Client{
		Id:                        client.Id,
		Active:                    client.Active,
		LockedUntil:               lockedUntil.Value(),
		Locked:                    client.Locked,
		CashAvailable:             client.CashAvailable,
		CashAmount:                client.CashAmount,
		CashReserved:              client.CashReserved,
		CryptoAvailable:           client.CryptoAvailable,
		CryptoAmount:              client.CryptoAmount,
		CryptoReserved:            client.CryptoReserved,
		OperationStopLoss:         client.OperationStopLoss,
		DayStopLoss:               client.DayStopLoss,
		MonthStopLoss:             client.MonthStopLoss,
		OperationAmountPercentage: client.OperationAmountPercentage,
		BuyOn:                     client.BuyOn,
		SellOn:                    client.SellOn,
		Symbols:                   client.Symbols,
		Summary:                   summaries,
	}
}

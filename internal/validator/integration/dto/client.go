package dto

import "github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"

// Client DynamoDB entity for crypto-robot.client repository
type Client struct {
	Id             string   `dynamodbav:"client_id"`
	Active         bool     `dynamodbav:"active"`
	LockedUntil    string   `dynamodbav:"locked_until"`
	Locked         bool     `dynamodbav:"locked"`
	CashAmount     float64  `dynamodbav:"cash_amount"`
	CashReserved   float64  `dynamodbav:"cash_reserved"`
	CryptoAmount   float64  `dynamodbav:"crypto_amount"`
	CryptoReserved float64  `dynamodbav:"crypto_reserved"`
	BuyOn          int      `dynamodbav:"buy_on"`
	SellOn         int      `dynamodbav:"sell_on"`
	Symbols        []string `dynamodbav:"symbols"`
}

// ClientDto creates a dto.Client from model.Client
func ClientDto(client *model.Client) *Client {
	return &Client{
		Id:     client.Id,
		Locked: client.Locked,
	}
}

// ToModel creates a model.Client from dto.Client
func (c Client) ToModel() *model.Client {
	return &model.Client{
		Id:     c.Id,
		Locked: c.Locked,
	}
}

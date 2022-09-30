package model

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/summary_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"time"
)

type Client struct {
	Id                        string
	Active                    bool
	LockedUntil               time.Time
	Locked                    bool
	CashAvailable             float64
	CashAmount                float64
	CashReserved              float64
	CryptoAvailable           float64
	CryptoAmount              float64
	CryptoReserved            float64
	OperationStopLoss         float64
	DayStopLoss               float64
	MonthStopLoss             float64
	OperationAmountPercentage float64
	BuyOn                     int
	SellOn                    int
	Symbols                   []string
	Summary                   []*Summary
}

// SetBalance will update client current balance, will take account of reserved values.
func (c *Client) SetBalance(balance *Balance) {
	c.CashAmount = balance.BrlBalance - c.CashReserved
	c.CryptoAmount = balance.CryptoBalance - c.CryptoReserved
}

// CreateOperation validates if client current values can operate, then creates a model.Operation and also updates
// reserved balance as necessary for the operation. Will return error in case of validation failure.
func (c *Client) CreateOperation(request *OperationRequest, coin *Coin) (*Operation, custom_error.BaseErrorAdapter) {
	timeUtils := time_utils.Time()
	for _, summary := range c.Summary {
		if summary.Type == summary_type.Day && timeUtils.IsToday(summary.Year, summary.Month, summary.Day) && summary.Profit < c.DayStopLoss*-1 {
			c.LockedUntil = timeUtils.Tomorrow()
			return nil, c.abort("Client day stop loss reached")
		}
		if summary.Type == summary_type.Month && timeUtils.IsThisMonth(summary.Year, summary.Month) && summary.Profit < c.MonthStopLoss*-1 {
			c.LockedUntil = timeUtils.NextMonth()
			return nil, c.abort("Client month stop loss reached")
		}
	}

	operation := NewOperation(c.OperationStopLoss)

	switch request.Operation {
	case operation_type.Buy:
		if coin.GetMinOperationValue(operation_type.Buy) > c.CashAmount || coin.GetMinOperationValue(operation_type.Buy) > c.CashAvailable {
			return nil, c.abort("Client does not have minimum cash amount")
		}

		expectedOperationAmount := c.CashAvailable * c.OperationAmountPercentage / 100

		if expectedOperationAmount > c.CashAmount {
			operation.Amount = c.CashAmount
			c.CashReserved += c.CashAmount
			c.CashAmount = 0
		} else {
			operation.Amount = expectedOperationAmount
			c.CashReserved += expectedOperationAmount
			c.CashAmount -= expectedOperationAmount
		}

		operation.Type = operation_type.Buy
		operation.Quote = symbol.Bitcoin
		operation.Base = symbol.Brl
	case operation_type.Sell:
		if coin.GetMinOperationValue(operation_type.Sell) > c.CryptoAmount || coin.GetMinOperationValue(operation_type.Sell) > c.CryptoAvailable {
			return nil, exceptions.NewValidationError("Client does not have minimum crypto amount")
		}

		expectedOperationAmount := c.CryptoAvailable * c.OperationAmountPercentage / 100

		if expectedOperationAmount > c.CryptoAmount {
			operation.Amount = c.CryptoAmount
			c.CryptoReserved += c.CryptoAmount
			c.CryptoAmount = 0
		} else {
			operation.Amount = expectedOperationAmount
			c.CryptoReserved += expectedOperationAmount
			c.CryptoAmount -= expectedOperationAmount
		}

		operation.Type = operation_type.Sell
		operation.Quote = symbol.Brl
		operation.Base = symbol.Bitcoin
	}

	return operation, nil
}

// Lock client
func (c *Client) Lock() {
	c.Locked = true
}

// Unlock client
func (c *Client) Unlock() {
	c.Locked = false
}

func (c *Client) abort(message string) custom_error.BaseErrorAdapter {
	return exceptions.NewValidationError(message)
}

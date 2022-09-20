package model

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

type Client struct {
	Id             string
	Active         bool
	LockedUntil    string
	Locked         bool
	CashAmount     float64
	CashReserved   float64
	CryptoAmount   float64
	CryptoReserved float64
	BuyOn          int
	SellOn         int
	Symbols        []string
}

// SetBalance will update client current balance, will take account of reserved values.
func (c *Client) SetBalance(balance *Balance) {

}

// CreateOperation validates if client current values can operate, then creates a model.Operation and also updates
// reserved balance as necessary for the operation. Will return error in case of validation failure.
func (c *Client) CreateOperation(request *OperationRequest, coin *Coin) (*Operation, custom_error.BaseErrorAdapter) {
	//TODO: validate client operation execution
	//TODO: create operation
	return nil, nil
}

// Lock client
func (c *Client) Lock() {
	c.Locked = true
}

// Unlock client
func (c *Client) Unlock() {
	c.Locked = false
}

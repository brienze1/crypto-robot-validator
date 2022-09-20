package model

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

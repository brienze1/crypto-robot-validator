package dto

// CoinResponse is the model for Biscoint GetCoin response.
type CoinResponse struct {
	Message string `json:"message"`
	Coin    Coin   `json:"data"`
}

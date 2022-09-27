package dto

// BalanceResponse is the model for Biscoint GetBalance response.
type BalanceResponse struct {
	Message string  `json:"message"`
	Balance Balance `json:"data"`
}

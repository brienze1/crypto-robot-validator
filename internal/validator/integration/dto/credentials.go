package dto

type Credentials struct {
	ClientId  string `json:"client_id" dynamodbav:"client_id"`
	ApiKey    string `json:"api_key" dynamodbav:"api_key"`
	ApiSecret string `json:"api_secret" dynamodbav:"api_secret"`
}

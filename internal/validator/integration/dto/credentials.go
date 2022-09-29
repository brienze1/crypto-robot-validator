package dto

type Credentials struct {
	ClientId  string `dynamodbav:"client_id"`
	ApiKey    string `dynamodbav:"api_key"`
	ApiSecret string `dynamodbav:"api_secret"`
}

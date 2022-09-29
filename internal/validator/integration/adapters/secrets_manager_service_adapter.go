package adapters

// SecretsManagerServiceAdapter is an adapter for secret manager service implementation.
type SecretsManagerServiceAdapter interface {
	GetSecret(secretName string, secretObject any) error
}

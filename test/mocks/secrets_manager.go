package mocks

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type secretsManager struct {
	GetSecretValueCounter int
	GetSecretValueError   error
	ReturnEmptyBinary     bool
	ReturnEmptyString     bool
	secrets               map[string][]byte
}

func SecretsManager() *secretsManager {
	return &secretsManager{}
}

func (s *secretsManager) GetSecretValue(_ context.Context, params *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (
	*secretsmanager.GetSecretValueOutput,
	error,
) {
	s.GetSecretValueCounter++

	if s.GetSecretValueError != nil {
		return nil, s.GetSecretValueError
	}

	secrets := s.secrets[*params.SecretId]
	encodedBinarySecretBytes := make([]byte, base64.StdEncoding.EncodedLen(len(secrets)))
	base64.StdEncoding.Encode(encodedBinarySecretBytes, secrets)
	var secretBinary []byte
	var secretString *string
	if !s.ReturnEmptyBinary {
		secretBinary = encodedBinarySecretBytes
	}
	if !s.ReturnEmptyString {
		secretString = aws.String(string(secrets))
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretBinary: secretBinary,
		SecretString: secretString,
	}, nil
}

func (s *secretsManager) SetSecret(key string, value any) {
	byteValue, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	s.secrets[key] = byteValue
}

func (s *secretsManager) Reset() {
	s.GetSecretValueCounter = 0
	s.GetSecretValueError = nil
	s.ReturnEmptyString = false
	s.ReturnEmptyBinary = false
	s.secrets = map[string][]byte{}
}

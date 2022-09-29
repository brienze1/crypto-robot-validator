package mocks

import "encoding/json"

type secretsManagerService struct {
	secrets map[string][]byte
}

func SecretsManagerService() *secretsManagerService {
	return &secretsManagerService{
		secrets: map[string][]byte{},
	}
}

func (s *secretsManagerService) GetSecret(secretName string, secretObject any) error {
	return json.Unmarshal(s.secrets[secretName], secretObject)
}

func (s *secretsManagerService) SetSecret(key string, value any) {
	byteValue, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	s.secrets[key] = byteValue
}

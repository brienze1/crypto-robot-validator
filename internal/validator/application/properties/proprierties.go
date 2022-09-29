package properties

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type properties struct {
	Profile                        string
	MinimumCryptoSellOperation     float64
	MinimumCryptoBuyOperation      float64
	BiscointUrl                    string
	SimulationUrl                  string
	BiscointGetCryptoPath          string
	BiscointGetBalancePath         string
	CryptoOperationTriggerTopicArn string
	Aws                            *aws
	Cache                          *cache
}

type cache struct {
	KeyTTL    time.Duration
	KeyPrefix string
}

type aws struct {
	Config         *awsConfig
	DynamoDB       *dynamoDB
	SecretsManager *secretsManager
}

type awsConfig struct {
	Region         string
	URL            string
	AccessKey      string
	AccessSecret   string
	Token          string
	OverrideConfig bool
}

type dynamoDB struct {
	ClientTableName      *string
	OperationTableName   *string
	CredentialsTableName *string
}

type secretsManager struct {
	CacheSecretName string
}

var once sync.Once

var propertiesInstance *properties

// Properties class is used to store and use env variables in runtime
func Properties() *properties {
	if propertiesInstance == nil {
		propertiesLoaded := loadProperties()
		once.Do(
			func() {
				propertiesInstance = propertiesLoaded
			})
	}

	return propertiesInstance
}

func loadProperties() *properties {
	profile := os.Getenv("PROFILE")
	minimumCryptoSellOperation := getDoubleEnvVariable("MINIMUM_CRYPTO_SELL_OPERATION")
	minimumCryptoBuyOperation := getDoubleEnvVariable("MINIMUM_CRYPTO_BUY_OPERATION")
	biscointUrl := os.Getenv("BISCOINT_CRYPTO_URL")
	biscointGetCryptoPath := os.Getenv("BISCOINT_CRYPTO_GET_CRYPTO_PATH")
	biscointGetBalancePath := os.Getenv("BISCOINT_CRYPTO_GET_BALANCE_PATH")
	cryptoOperationTriggerTopicArn := os.Getenv("AWS_SNS_TOPIC_ARN_CRYPTO_OPERATIONS")
	awsRegion := os.Getenv("AWS_REGION")
	awsURL := os.Getenv("AWS_URL")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsAccessSecret := os.Getenv("AWS_ACCESS_SECRET")
	awsAccessToken := os.Getenv("AWS_ACCESS_TOKEN")
	awsOverrideConfig := getBoolEnvVariable("AWS_OVERRIDE_CONFIG")
	clientTableName := os.Getenv("AWS_DYNAMODB_CLIENT_TABLE_NAME")
	operationTableName := os.Getenv("AWS_DYNAMODB_OPERATION_TABLE_NAME")
	credentialsTableName := os.Getenv("AWS_DYNAMODB_CREDENTIALS_TABLE_NAME")
	cacheSecretName := os.Getenv("AWS_SECRETS_MANAGER_CACHE_SECRET_NAME")
	cacheKeyTTL := getIntEnvVariable("CACHE_KEY_TTL_SECONDS")
	cacheKeyPrefix := os.Getenv("CACHE_KEY_PREFIX")

	return &properties{
		Profile:                        profile,
		MinimumCryptoSellOperation:     minimumCryptoSellOperation,
		MinimumCryptoBuyOperation:      minimumCryptoBuyOperation,
		SimulationUrl:                  biscointUrl,
		BiscointUrl:                    biscointUrl,
		BiscointGetCryptoPath:          biscointGetCryptoPath,
		BiscointGetBalancePath:         biscointGetBalancePath,
		CryptoOperationTriggerTopicArn: cryptoOperationTriggerTopicArn,
		Aws: &aws{
			Config: &awsConfig{
				Region:         awsRegion,
				URL:            awsURL,
				AccessKey:      awsAccessKey,
				AccessSecret:   awsAccessSecret,
				Token:          awsAccessToken,
				OverrideConfig: awsOverrideConfig,
			},
			DynamoDB: &dynamoDB{
				ClientTableName:      &clientTableName,
				OperationTableName:   &operationTableName,
				CredentialsTableName: &credentialsTableName,
			},
			SecretsManager: &secretsManager{
				CacheSecretName: cacheSecretName,
			},
		},
		Cache: &cache{
			KeyTTL:    time.Duration(cacheKeyTTL) * time.Second,
			KeyPrefix: cacheKeyPrefix,
		},
	}
}

func getDoubleEnvVariable(key string) float64 {
	value, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		panic(err.Error() + ". Failed to load property \"" + key + "\" from environment")
	}

	return value
}

func getBoolEnvVariable(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}

	return value
}

func getIntEnvVariable(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err.Error() + ". Failed to load property \"" + key + "\" from environment")
	}

	return value
}

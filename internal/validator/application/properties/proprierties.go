package properties

import (
	"os"
	"strconv"
	"sync"
)

type properties struct {
	Profile                           string
	MinimumCryptoSellOperation        float64
	MinimumCryptoBuyOperation         float64
	BinanceCryptoSymbolPriceTickerUrl string
	CryptoOperationTriggerTopicArn    string
	Aws                               *aws
}

type aws struct {
	Config   *awsConfig
	DynamoDB *dynamoDB
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
	ClientTableName string
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
	binanceCryptoSymbolPriceTickerUrl := os.Getenv("BINANCE_CRYPTO_SYMBOL_PRICE_TICKER_URL")
	cryptoOperationTriggerTopicArn := os.Getenv("AWS_SNS_TOPIC_ARN_CRYPTO_OPERATIONS")
	awsRegion := os.Getenv("AWS_REGION")
	awsURL := os.Getenv("AWS_URL")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsAccessSecret := os.Getenv("AWS_ACCESS_SECRET")
	awsAccessToken := os.Getenv("AWS_ACCESS_TOKEN")
	awsOverrideConfig := getBoolEnvVariable("AWS_OVERRIDE_CONFIG")
	clientTableName := os.Getenv("AWS_DYNAMODB_CLIENT_TABLE_NAME")

	return &properties{
		Profile:                           profile,
		MinimumCryptoSellOperation:        minimumCryptoSellOperation,
		MinimumCryptoBuyOperation:         minimumCryptoBuyOperation,
		BinanceCryptoSymbolPriceTickerUrl: binanceCryptoSymbolPriceTickerUrl,
		CryptoOperationTriggerTopicArn:    cryptoOperationTriggerTopicArn,
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
				ClientTableName: clientTableName,
			},
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

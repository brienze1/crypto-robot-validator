module github.com/brienze1/crypto-robot-validator

go 1.19

require (
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go-v2 v1.16.15
	github.com/aws/aws-sdk-go-v2/config v1.17.6
	github.com/aws/aws-sdk-go-v2/credentials v1.12.19
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.9.18
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.17.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.16.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.18.0
	github.com/cucumber/godog v0.12.5
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.18 // indirect
	github.com/aws/smithy-go v1.13.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-memdb v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gogo/protobuf v1.2.1 => github.com/gogo/protobuf v1.3.2
	github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.50
	github.com/prometheus/client_golang v0.9.3 => github.com/prometheus/client_golang v1.13.0
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.8
)

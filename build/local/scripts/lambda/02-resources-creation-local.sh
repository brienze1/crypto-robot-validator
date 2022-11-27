#!/bin/bash


echo "-----------------Script-02----------------- [validator]"

echo "########### Create secrets manager for cache values ###########"
aws secretsmanager create-secret --name crypto_robot.secrets.cache --secret-string '{"redis_address":"localhost:6379","redis_password":"password123","redis_user":"default"}' --endpoint-url http://localstack:4566

echo "########### Create secrets manager for encryption values ###########"
aws secretsmanager create-secret --name crypto_robot.secrets.encryption --secret-string '{"encryption_key":"9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@"}' --endpoint-url http://localstack:4566

echo "########### Copy the lambda function to the S3 bucket ###########"
aws s3 cp /lambda-files/crypto-robot-validator.zip s3://lambda-functions --endpoint-url http://localstack:4566

echo "########### Cloudformation Start ###########"
aws cloudformation deploy \
 --stack-name crypto-robot-validator \
 --template-file "/cloudformation/crypto-robot-validator.yaml" \
 --endpoint-url http://localstack:4566

#echo "########### Creating DLQ SQS ###########"
#aws sqs create-queue --queue-name cryptoValidatorQueueDLQ --endpoint-url http://localstack:4566
#
#echo "########### Creating SQS ###########"
#aws sqs create-queue \
#--queue-name cryptoValidatorQueue \
#--attributes '{
#    "RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:sa-east-1:000000000000:cryptoValidatorQueueDLQ\",\"maxReceiveCount\":\"3\"}",
#    "MessageRetentionPeriod": "259200",
#    "VisibilityTimeout": "90"
#}' \
#--endpoint-url http://localstack:4566
#
#echo "########### Listing SQS ###########"
#aws sqs list-queues --endpoint-url http://localstack:4566
#
#echo "########### Subscribing SQS to SNS ###########"
#aws sns subscribe \
#--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoOperationTriggerTopic \
#--protocol sqs \
#--notification-endpoint "http://localhost:4566/000000000000/cryptoValidatorQueue" \
#--endpoint-url http://localstack:4566
#
#echo "########### Listing SNS Subscriptions ###########"
#aws sns list-subscriptions --endpoint-url http://localstack:4566
#
#echo "########### Creating SNS ###########"
#aws sns create-topic --name cryptoOperationExecutorTopic --endpoint-url http://localstack:4566
#
#echo "########### Listing SNS ###########"
#aws sns list-topics --endpoint-url http://localstack:4566
#
#echo "########### Create secrets manager for cache values ###########"
#aws secretsmanager create-secret --name crypto_robot.secrets.cache --secret-string '{"redis_address":"localhost:6379","redis_password":"password123","redis_user":"default"}' --endpoint-url http://localstack:4566
#
#echo "########### Create secrets manager for encryption values ###########"
#aws secretsmanager create-secret --name crypto_robot.secrets.encryption --secret-string '{"encryption_key":"9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@"}' --endpoint-url http://localstack:4566
#
#echo "########### Copy the lambda function to the S3 bucket ###########"
#aws s3 cp /lambda-files/crypto-robot-validator.zip s3://lambda-functions --endpoint-url http://localstack:4566
#
#echo "########### Create the lambda validatorLambda ###########"
#aws lambda create-function \
#  --endpoint-url http://localstack:4566 \
#  --function-name validatorLambda \
#  --role arn:aws:iam::000000000000:role/admin-role \
#  --code S3Bucket=lambda-functions,S3Key=crypto-robot-validator.zip \
#  --handler ./validator \
#  --runtime go1.x \
#  --description "SQS Lambda handler for crypto-robot-validator." \
#  --timeout 60 \
#  --memory-size 128 \
#  --environment "Variables={VALIDATOR_ENV=localstack}"
#
#echo "########### Map the cryptoValidatorQueue to the validatorLambda lambda function ###########"
#aws lambda create-event-source-mapping \
#  --function-name validatorLambda \
#  --batch-size 1 \
#  --event-source-arn "arn:aws:sqs:sa-east-1:000000000000:cryptoValidatorQueue" \
#  --endpoint-url http://localstack:4566
#
## Only need this for local testing
#echo "########### Creating DynamoDB 'crypto_robot.clients' table ###########"
#aws dynamodb create-table \
#--table-name crypto_robot.clients  \
#--attribute-definitions AttributeName=client_id,AttributeType=S \
#--key-schema AttributeName=client_id,KeyType=HASH \
#--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
#--endpoint-url=http://localstack:4566
#
#echo "########### Creating DynamoDB 'crypto_robot.operations' table ###########"
#aws dynamodb create-table \
#--table-name crypto_robot.operations  \
#--attribute-definitions AttributeName=operation_id,AttributeType=S \
#--key-schema AttributeName=operation_id,KeyType=HASH \
#--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
#--endpoint-url=http://localstack:4566
#
#echo "########### Creating DynamoDB 'crypto_robot.credentials' table ###########"
#aws dynamodb create-table \
#--table-name crypto_robot.credentials  \
#--attribute-definitions AttributeName=client_id,AttributeType=S \
#--key-schema AttributeName=client_id,KeyType=HASH \
#--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
#--endpoint-url=http://localstack:4566
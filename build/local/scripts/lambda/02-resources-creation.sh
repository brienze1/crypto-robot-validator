#!/bin/bash


echo "-----------------Script-02----------------- [operation-hub]"

echo "########### Creating SNS ###########"
aws sns create-topic --name cryptoAnalysisSummaryTopic --endpoint-url http://localstack:4566

echo "########### Creating SNS ###########"
aws sns create-topic --name cryptoOperationTriggerTopic --endpoint-url http://localstack:4566

echo "########### Listing SNS ###########"
aws sns list-topics --endpoint-url http://localstack:4566

echo "########### Creating SQS ###########"
aws sqs create-queue --queue-name cryptoOperationHubQueue --endpoint-url http://localstack:4566

echo "########### Listing SQS ###########"
aws sqs list-queues --endpoint-url http://localstack:4566

echo "########### Subscribing SQS to SNS ###########"
aws sns subscribe \
--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoAnalysisSummaryTopic \
--protocol sqs \
--notification-endpoint "http://localhost:4566/000000000000/cryptoOperationHubQueue" \
--endpoint-url http://localstack:4566

echo "########### Listing SNS Subscriptions ###########"
aws sns list-subscriptions --endpoint-url http://localstack:4566

echo "########### Create secrets manager ###########"
aws secretsmanager create-secret --name cryptoRobotOperationHubSecret --secret-string '{"host":"localhost","port":5432,"user":"postgres","password":"postgres","db_name":"crypto_robot"}' --endpoint-url http://localstack:4566

echo "########### Copy the lambda function to the S3 bucket ###########"
aws s3 cp /lambda-files/crypto-robot-operation-hub.zip s3://lambda-functions --endpoint-url http://localstack:4566

echo "########### Create the lambda operationHubLambda ###########"
aws lambda create-function \
  --endpoint-url http://localstack:4566 \
  --function-name operationHubLambda \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=crypto-robot-operation-hub.zip \
  --handler ./operation-hub \
  --runtime go1.x \
  --description "SQS Lambda handler for test sqs." \
  --timeout 60 \
  --memory-size 128 \
  --environment "Variables={OPERATION_HUB_ENV=localstack}"

echo "########### Map the cryptoOperationHubQueue to the operationHubLambda lambda function ###########"
aws lambda create-event-source-mapping \
  --function-name operationHubLambda \
  --batch-size 1 \
  --event-source-arn "arn:aws:sqs:sa-east-1:000000000000:cryptoOperationHubQueue" \
  --endpoint-url http://localstack:4566

echo "########### Creating DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb create-table \
--table-name crypto_robot.clients  \
--attribute-definitions AttributeName=client_id,AttributeType=S \
--key-schema AttributeName=client_id,KeyType=HASH \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
--endpoint-url=http://localstack:4566


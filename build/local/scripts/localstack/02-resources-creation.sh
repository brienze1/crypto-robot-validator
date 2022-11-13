#!/bin/bash

echo "-----------------Script-02----------------- [localstack]"


echo "########### Make S3 bucket for lambdas ###########"
aws s3 mb s3://lambda-functions --endpoint-url http://localhost:4566

echo "########### Create Admin IAM Role ###########"
aws iam create-role --role-name admin-role --path / --assume-role-policy-document file:./admin-policy.json --endpoint-url http://localhost:4566

echo "########### Creating SNS ###########"
aws sns create-topic --name cryptoOperationTriggerTopic --endpoint-url http://localhost:4566

echo "########### Listing SNS ###########"
aws sns list-topics --endpoint-url http://localhost:4566

echo "########### Creating DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb create-table \
--table-name crypto_robot.clients  \
--attribute-definitions AttributeName=client_id,AttributeType=S \
--key-schema AttributeName=client_id,KeyType=HASH \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
--endpoint-url=http://localstack:4566
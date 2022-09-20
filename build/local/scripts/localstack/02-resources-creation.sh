#!/bin/bash

echo "-----------------Script-02----------------- [localstack]"


echo "########### Make S3 bucket for lambdas ###########"
aws s3 mb s3://lambda-functions --endpoint-url http://localhost:4566

echo "########### Create Admin IAM Role ###########"
aws iam create-role --role-name admin-role --path / --assume-role-policy-document file:./admin-policy.json --endpoint-url http://localhost:4566
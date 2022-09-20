#!/bin/bash


echo "-----------------Script-01----------------- [operation-hub]"

echo "########### Check if localstack is up ###########"
until curl http://localstack:4566/health --silent; do
  echo "Localstack not up yet"
  sleep 1
done

echo "########### Check if s3 lambda bucket was created ###########"
until aws s3 ls s3://lambda-functions --endpoint-url=http://localstack:4566; do
  echo "S3 bucket \"lambda-functions\" not created yet"
  sleep 1
done

echo "########### Check if admin IAM role was created ###########"
until aws iam get-role --role-name "admin-role" --endpoint-url=http://localstack:4566; do
  echo "IAM role \"admin-role\" not created yet"
  sleep 1
done

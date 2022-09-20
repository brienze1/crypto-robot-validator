package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/brienze1/crypto-robot-validator/internal/validator"
)

func main() {
	lambda.Start(validator.Main().Handle)
}

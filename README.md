<div id="top"></div>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/brienze1/crypto-robot-validator/blob/main/LICENSE)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/brienze1/crypto-robot-validator)
![Build](https://img.shields.io/github/workflow/status/brienze1/crypto-robot-validator/Build?label=Build)
[![Coverage Status](https://coveralls.io/repos/github/brienze1/crypto-robot-validator/badge.svg?branch=main)](https://coveralls.io/github/brienze1/crypto-robot-validator?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/brienze1/crypto-robot-validator)](https://goreportcard.com/report/github.com/brienze1/crypto-robot-validator)
[![Golang](https://img.shields.io/github/go-mod/go-version/brienze1/crypto-robot-validator)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/brienze1/crypto-robot-validator.svg)](https://pkg.go.dev/github.com/brienze1/crypto-robot-validator)

# Crypto Data Validator

1. [About the Project](#about-the-project)
    1. [Input](#input)
    2. [Output](#output)
    3. [Persistence](#persistence)
        1. [Client DB](#client-db)
            1. [Schema](#schema)
            2. [Operation](#operation)
            3. [Query](#query)
    4. [Rules](#rules)
    5. [Built With](#built-with)
        1. [Dependencies](#dependencies)
        2. [Compiler Dependencies](#compiler-dependencies)
        3. [Test Dependencies](#test-dependencies)
    6. [Roadmap](#roadmap)
2. [Getting Started](#getting-started)
    1. [Prerequisites](#prerequisites)
    2. [Installation](#installation)
    3. [Requirements](#requirements)
        1. [Deploying Local Infrastructure](#deploying-local-infrastructure)
    4. [Usage](#usage)
        1. [Manual Input](#manual-input)
        2. [Docker Input](#docker-input)
    5. [Testing](#testing)
3. [About Me](#about-me)

## About the Project

The objective of this project is to receive a BUY or SELL event and trigger operations for active users.

### Input

The input should be received as an SNS message sent through an SQS subscription. This message will trigger the lambda
handler to perform the service. The indicators received should be used to predict marked behaviour and trigger BUY or
SELL operations.

Data could also be analysed for specific interval indicators if needed in the future (but the current behaviour does not
use that data).

Analysis summary indicators:

- STRONG_BUY
- BUY
- NEUTRAL
- SELL
- STRONG_SELL

Example of how the data received should look like:

```json
{
  "summary": "BUY",
  "timestamp": "20-07-2022 02:18:10",
  "analysed_data": [
    {
      "interval": "0NE_MINUTE",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "BUY",
      "analysis": [
        {
          "metric": "SIMPLE_MOVING_AVERAGE",
          "indicator": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "FIVE_MINUTES",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "FIFTEEN_MINUTES",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "THIRTY_MINUTES",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "ONE_HOUR",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "SIX_HOURS",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    },
    {
      "interval": "ONE_DAY",
      "timestamp": "20-07-2022 02:18:10",
      "summary": "STRONG_BUY",
      "analysis": [
        {
          "indicator": "SIMPLE_MOVING_AVERAGE",
          "summary": "BUY",
          "score": {
            "buy": 4,
            "sell": 2
          }
        },
        {
          "indicator": "EXPONENTIAL_MOVING_AVERAGE",
          "summary": "NEUTRAL",
          "score": {
            "buy": 3,
            "sell": 3
          }
        }
      ]
    }
  ]
}
```

### Output

Since this is an async application there is no output to be returned, but operation events are generated from the data
received.
Operation events should have information needed to track client who will trigger the operation, the operation type,
crypto being traded, operation event time sent and available amount (if operation type is BUY amount should be in cash,
otherwise if operation type is SELL, amount should be used like cryptocurrency amount).

Example of how the line should look like:

```
{
    "client_id": "asdasd-asdasd-asdasd-asdasd",
    "operation": "BUY",
    "symbol": "BTC",
    "start_time": "2007-12-03 10:15:30",
}
```

### Persistence

#### Client DB

Client DB is the database that contains the client information and configuration needed to trigger the operations.
For this DB, DynamoDB was chosen because it's cheaper.

##### Schema

[//]: # (TODO fix schema)

```json
{
  "id": "uuid",
  "active": true,
  "locked_until": "2022-09-17T12:05:07.45066-03:00",
  "locked": false,
  "cash_amount": 100,
  "cash_reserved": 0.00,
  "crypto_amount": 0.0000312,
  "crypto_reserved": 0.0,
  "symbols": [
    "BTC",
    "SOL"
  ],
  "buy_on": "STRONG_BUY",
  "sell_on": "SELL",
  "ops_timeout_seconds": 60,
  "operation_stop_loss": 50.00,
  "day_stop_loss": 500.00,
  "month_stop_loss": 500.00,
  "summary": [
    {
      "type": "MONTH",
      "day": 1,
      "month": 8,
      "year": 2022,
      "amount_sold": 23000.42,
      "amount_bought": 37123.42,
      "profit": 1032.32,
      "crypto": [
        {
          "symbol": "BTC",
          "average_buy_value": 230020.42,
          "average_sell_value": 235020.42,
          "amount_sold": 0.00231,
          "amount_bought": 0.00431,
          "profit": -53.00
        }
      ]
    },
    {
      "type": "DAY",
      "day": 14,
      "month": 8,
      "year": 2022,
      "amount_sold": 23000.42,
      "amount_bought": 37123.42,
      "profit": -53.00,
      "crypto": [
        {
          "symbol": "BTC",
          "average_buy_value": 230020.42,
          "average_sell_value": 235020.42,
          "amount_sold": 0.00231,
          "amount_bought": 0.00431,
          "profit": -53.00
        }
      ]
    }
  ]
}
```

##### Operation

This application supports the following operations to the Client DB:

- Read ops:
    - Used to read clients available for the current operation

##### Query

This is the query used to get clients from DB:

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("active").Equal(expression.Value(config.Active)),
    expression.Name("locked").Equal(expression.Value(config.Locked)),
    expression.Name("locked_until").LessThanEqual(expression.Value(d.timeSource.Now())),
    expression.Name("cash_amount").GreaterThanEqual(expression.Value(config.MinimumCash)),
    expression.Name("crypto_amount").GreaterThanEqual(expression.Value(config.MinimumCrypto)),
    expression.Name("sell_on").LessThanEqual(expression.Value(config.SellWeight.Value())),
    expression.Name("buy_on").LessThanEqual(expression.Value(config.BuyWeight.Value())),
    expression.Name("symbols").Contains(config.Symbol.Name()),
    ),
    ).Build()
```

### Rules

Here are some rules that need to be implemented in this application.

Implemented:

- Client must be active
- Client must not be locked
- Current date must be greater than locked_until value
- Client must have enough cash to buy minimum allowed amount of crypto
- Client must have enough crypto to sell minimum allowed amount
- Client must have the coin symbol selected inside `config.symbols` variable to operate it
- Buy operations should be triggered when the summary received is equal or less restricting than the `config.buy_on`
  value.
    - For example if the config value is equal to `BUY` and a `STRONG_BUY` analysis was received, the operation should
      be allowed, and the opposite should be denied.
- Sell operations should be triggered when the summary received is equal or less restricting than the `config.sell_on`
  value.
    - For example if the config value is equal to `SELL` and a `STRONG_SELL` analysis was received, the operation should
      be allowed, and the opposite should be denied.

Not implemented (this will be done on another service):

- Operations should not be triggered if `daily_summary.proffit` has a negative value of more than or equal to
  the `config.day_stop_loss` value.
    - `daily_summary.day` value should be checked to see if current day has changed, in this case, the values
      should be updated to start a new day.
- Operations should not be triggered if `monthly_summary.proffit` has a negative value of more than or equal to
  the `config.month_stop_loss` value.
    - `monthly_summary.month` value should be checked to see if current month has changed, in this case, the values
      should be updated to start a new month.

### Built With

This application is build with Golang, code is build using a Dockerfile every deployment into the main branch in GitHub
using GitHub actions. Local environment is created using localstack for testing purposes using
[crypto-robot-localstack](https://github.com/brienze1/crypto-robot-localstack).

#### Dependencies

- [aws/aws-lambda-go](https://github.com/aws/aws-lambda-go): Used in Lambda Handler integration
- [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2): Used in SNS and DynamoDB integration
- [google/uuid](https://github.com/google/uuid): Used to generate uuids
- [joho/godotenv](https://github.com/joho/godotenv): Used to map .env variables

#### Compiler Dependencies

- [golangci/golangci-lint](https://github.com/golangci/golangci-lint): Used to enforce coding practices

#### Test Dependencies

- [cucumber/godog](https://github.com/cucumber/godog): Used to run integration tests
- [stretchr/testify](https://github.com/stretchr/testify): Used to perform test assertions

### Roadmap

- [x] Implement Behaviour tests (BDD)
- [x] Implement Unit tests
- [x] Implement application logic
- [x] Create Dockerfile
- [x] Create Docker compose for local infrastructure
- [x] Document everything in Readme
- [X] Change to use DynamoDB instead of Postgres (worse performance but, less expensive)
- [x] Use secret manager to get DB password

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

### Prerequisites

- Install Golang

    - Windows/MacOS/Linux
        - [Manual](https://go.dev/dl/)
    - macOS
        - [Homebrew](https://docs.brew.sh/Installation)
          ```bash
          brew install go
          ```
    - Linux
        - Via terminal
          ```bash
          sudo add-apt-repository ppa:longsleep/golang-backports
          sudo apt update
          sudo apt install golang-go
          ```

- Install Docker
    - [Windows/macOS/Linux/WSL](https://www.docker.com/get-started/)

### Installation

- Run the following to install project dependencies:
    - Windows/MacOS/Linux/WSL
      ```bash
      go mod download
      ```

- Run the following to compile the project and generate executable:
    - Windows/MacOS/Linux/WSL
      ```bash
      go build -o bin/validator cmd/validator/main.go
      ```

Note: the binary generated will be available at `./bin` folder.

### Requirements

To run the application locally, first a local infrastructure needs to be deployed

#### Deploying Local Infrastructure

This requires [docker](#prerequisites) to be installed. Localstack will deploy aws local integration and create the
topic used by this application to send the events.

Obs: Make sure Docker is running before.

- Start the required infrastructure via localstack using docker compose command:

    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose -f ./build/local/docker-compose.yml up
      ```

- To stop localstack:
    - Windows/macOS/Linux/WSL
      ```bash
      docker-compose -f ./build/local/docker-compose.yml down
      ```

### Usage

#### Manual Input

- Start the compiled application locally:
    - Windows/macOS/Linux/WSL
      ```bash
      go run cmd/local/main_local.go
      ```
- To stop the application just press Ctrl+C

#### Docker Input

- In case you want to use a Docker container to run the application first you need to build the Docker image from
  Dockerfile:
    - Windows/macOS/Linux/WSL
      ```bash
      docker build -t crypto-robot-validator .
      ```

- And then run the new created image:
    - Windows/macOS/Linux/WSL
      ```sh
      docker run --network="host" -d -it crypto-robot-validator bash \
      -c "VALIDATOR_ENV=localstack go run ./cmd/local/main_local.go"
      ```

### Testing

- To run the unit tests:
    - Windows/macOS/Linux/WSL
      ```bash
      go test ./test/unit/...
      ```

- To run the integration tests:
    - Windows/macOS/Linux/WSL
      ```bash
      go test ./test/integrated/...
      ```
      
- To run the all tests:
    - Windows/macOS/Linux/WSL
      ```bash
      go test ./...
      ```

<p align="right">(<a href="#top">back to top</a>)</p>

## About me

Hello! :)

My name is Luis Brienze, and I'm a Software Engineer.

I focus primarily on software development, but I'm also good at system architecture, mentoring other developers,
etc... I've been in the IT industry for 4+ years, during this time I worked for companies like Itaú, Dock, Imagine
Learning and
EPAM.

I graduated from UNESP studying Automation and Control Engineering in 2022, and I also took multiple courses on Udemy
and Alura.

My main stack is Java, but I'm also pretty good working with Kotlin and Typescript (both server side).
I have quite a good knowledge of AWS Cloud, and I'm also very conformable working with Docker.

During my career, while working with QA's, I've also gained a lot of valuable experience with testing applications in
general from unit/integrated
testing using TDD and BDD, to performance testing apps with JMeter for example.

If you want to talk to me, please fell free to reach me anytime at [LinkedIn](https://www.linkedin.com/in/luisbrienze/)
or [e-mail](mailto:lfbrienze@gmail.com?subject=[GitHUB]%20Crypto%20robot%20validator).

<p align="right">(<a href="#top">back to top</a>)</p>
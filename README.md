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
            1. [Schema](#client-db-schema)
            2. [Operation](#client-db-operation)
            3. [Query](#client-db-query)
        2. [Operation DB](#operation-db)
            1. [Schema](#operation-db-schema)
            2. [Operation](#operation-db-operation)
            3. [Query](#operation-db-query)
        3. [Lock DB](#lock-db)
            1. [Schema](#lock-db-schema)
            2. [Operation](#lock-db-operation)
            3. [Query](#lock-db-query)
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

The objective of this project is to receive an operation event for a single client, validate it and approve operation
execution.

### Input

The input should be received as an SNS message sent through an SQS subscription. This message will trigger the lambda
handler to perform the service. Operation events should have information needed to track client who will execute the
operation, the operation type, crypto being traded, operation event time sent and available amount (if operation type is
BUY amount should be in cash, otherwise if operation type is SELL, amount should be used like cryptocurrency amount).

Also, this step should validate user stop losses configuration and may lock the client until top loss is available
again.

Example of how the data received should look like:

```json
{
  "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
  "operation": "BUY",
  "symbol": "BTC",
  "start_time": "2022-09-17T12:05:07.45066-03:00"
}
```

### Output

Since this is an async application there is no output to be returned, but operation events are generated from the data
received.
Operation events should have information needed to get the operation from its db, the operation id.

Example of how the line should look like:

```json
{
  "operation_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"
}
```

### Persistence

#### Client DB

Client DB is the database that contains the client information and configuration needed to trigger the operations.

##### Client DB Schema

[//]: # (TODO fix schema)

```json
{
  "id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
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

##### Client DB Operation

This application supports the following operations to the Client DB:

- Read ops:
    - Used to find clients using client_id

- Write ops:
    - Used to lock clients using client_id

##### Client DB Query

This is the query used to get clients from DB:

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("client_id").Equal(expression.Value(client_id))),
    ),
    ).Build()
```

This is the query used to update (lock/unlock) clients:

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("client_id").Equal(expression.Value(client_id))),
    ),
    ).Build()
```

#### Operation DB

Operation DB is the database that contains the created operation's information.

##### Operation DB Schema

Operation statuses:

- CREATED
- PENDING
- COMPLETED
- ERROR

```json
{
  "id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
  "status": "COMPLETED",
  "created_at": "2022-09-17T12:05:07.45066-03:00",
  "expires_at": "2022-09-17T12:05:07.45066-03:00",
  "completed_at": "2022-09-17T12:05:07.45066-03:00",
  "locked": false,
  "quote": "BTC",
  "base": "BRL",
  "type": "BUY",
  "amount": 100.00,
  "stop_loss": 50.00,
  "profit": 1.0,
  "transactions": [
    {
      "type": "BUY",
      "unitary_value": 100000.00,
      "quote_amount": 0.001,
      "base_amount": 100.00,
      "created_at": "2022-09-17T12:05:07.45066-03:00",
      "expires_at": "2022-09-17T12:05:07.45066-03:00",
      "confirmed_at": "2022-09-17T12:05:07.45066-03:00"
    },
    {
      "type": "SELL",
      "unitary_value": 101000.00,
      "quote_amount": 0.001,
      "base_amount": 101.00,
      "created_at": "2022-09-17T12:05:07.45066-03:00",
      "expires_at": "2022-09-17T12:05:07.45066-03:00",
      "confirmed_at": "2022-09-17T12:05:07.45066-03:00"
    }
  ]
}
```

##### Operation DB Operation

This application supports the following operations to the Client DB:

- Write ops:
    - Used to create new operations

##### Operation DB Query

This is the query used to create operations in DB:

[//]: # (TODO create query)

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("client_id").Equal(expression.Value(client_id))),
    ),
    ).Build()
```

#### Lock DB

Lock DB is the database that contains the client_id's locked during execution.

OBS: Redis is used for this DB.

##### Lock DB Schema

This DB uses key value to store validator client_id's locked

```json
{
  "VALIDATOR_LOCK_{client_id}": "{client_id}"
}
```

##### Lock DB Operation

This application supports the following operations to the Lock DB:

- Read ops:
    - Used to find locked client_ids

- Write ops:
    - Used to lock client_ids

##### Lock DB Query

This is the query used to get locked client_id's from DB:

[//]: # (TODO create query)

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("client_id").Equal(expression.Value(client_id))),
    ),
    ).Build()
```

This is the query used to update (lock/unlock) client_id's from DB:

[//]: # (TODO create query)

```gotemplate
    expr, _ := expression.NewBuilder().WithFilter(
    expression.And(
    expression.Name("client_id").Equal(expression.Value(client_id))),
    ),
    ).Build()
```

### Rules

Here are some rules that need to be implemented in this application.

Not Implemented:

Client validations:

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
- Operations should not be triggered if `daily_summary.proffit` has a negative value of more than or equal to
  the `config.day_stop_loss` value.
    - `daily_summary.day` value should be checked to see if current day has changed, in this case, the values
      should be updated to start a new day.
- Operations should not be triggered if `monthly_summary.proffit` has a negative value of more than or equal to
  the `config.month_stop_loss` value.
    - `monthly_summary.month` value should be checked to see if current month has changed, in this case, the values
      should be updated to start a new month.

Lock:

- Ids received should be locked on Redis for execution and unlocked after, even if error occurred.
- Clients should be locked on DynamoDB for execution and unlocked after, if error occurred after DynamoDB lock, clients
  should be
  unlocked.
- If client fails validation `locked_until` value could be set on DynamoDB to lock for an extended amount of time (stop
  loss block for example)

Operations:

- Operation should be created with status `CREATED` and it's id should be sent to the SNS topic for later execution.
- Operation amount should be created using client configuration and Biscoint current unitary value.

Biscoint:

- Client balance should be validated from Biscoint and updated in DynamoDB clients DB.

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

- [ ] Implement Behaviour tests (BDD)
- [ ] Implement Unit tests
- [ ] Implement application logic
- [x] Create Dockerfile
- [ ] Create Docker compose for local infrastructure
- [ ] Document everything in Readme

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

My main stack is Java, but I'm also pretty good working with Kotlin, Typescript and Go (backend only). I have quite a
good knowledge of AWS Cloud, and I'm also very conformable working with Docker.

Also, I have experience working with relational (PostgreSQL, Microsoft SQL Server, MySQL, ...) and non-relational (
DynamoDB, Redis, Cassandra, ...) databases.

During my career, while working with QA's, I've also gained a lot of valuable experience with testing applications in
general from unit/integrated testing using TDD and BDD, to performance testing apps with JMeter for example.

If you want to talk to me, please fell free to reach me anytime at [LinkedIn](https://www.linkedin.com/in/luisbrienze/)
or [e-mail](mailto:lfbrienze@gmail.com?subject=[GitHUB]%20Crypto%20robot%20validator).

<p align="right">(<a href="#top">back to top</a>)</p>

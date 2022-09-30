#noinspection CucumberUndefinedStep
Feature: Validate operation
  In order to validate and trigger operations
  The client must be available and its credentials should be available as well
  The client should not be locked on redis
  The client should not be locked on dynamodb
  The client daily loss should be less than the configured daily stop loss
  The client monthly loss should be less than the configured monthly stop loss
  The client must have enough balance on biscoint
  An operation request must be received via sns

  Background:
    Given test env variables were loaded
    And dynamoDB is up
    And redis is up
    And biscoint api is up
    And sns service is up
    And secrets manager service is up

  Scenario: Validate operation request for one client with success
    Given there is a client available on DynamoDB with client id "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"
    And client available "brl" balance is 10000.00
    And client available "btc" balance is 0.1
    And client reserved "brl" balance is 1000.00
    And client reserved "btc" balance is 0.0001
    And client "brl" balance is 10000.00 on biscoint
    And client "btc" balance is 0.0001 on biscoint
    And crypto current "buy" value is 100000.00 on biscoint
    And crypto current "sell" value is 99000.00 on biscoint
    And the following credentials available for client id "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"
      """
      {
        "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
        "api_key": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
        "api_secret": "a7aca6d4f67519fbb4dc65b159b4e9526b069a2cb5f515d4690bce05ba81e6e5967f477e0ce3affa7c80843f3efed1cee9b0c062"
      }
      """
    When the following message is received
      """
      {
        "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
        "operation": "BUY",
        "symbol": "BTC",
        "start_time": "2022-09-17T12:05:07.45066-03:00"
      }
      """
    Then there should be 1 messages sent via sns
    And process should exit with 0


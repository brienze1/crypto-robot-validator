#noinspection CucumberUndefinedStep
Feature: Trigger operations
  In order to trigger operations
  There must be at least one client available
  I must receive a summary analysis via sqs

  Scenario: Trigger operation for one client
    Given test env variables were loaded
    And dynamoDB is "up"
    And binance api is "up"
    And sns service is "up"
    And I receive message with summary equals "STRONG_BUY"
    And there are 1 clients available in DB
    When handler is triggered
    Then there should be 1 messages sent via sns
    And sns messages payload should have all client_id's got from clients table
    And sns messages payload symbol should be equal "STRONG_BUY"
    And sns messages payload operation should be equal "BUY"
    And process should exit with 0

  Scenario: Trigger operation for six clients
    Given test env variables were loaded
    And dynamoDB is "up"
    And binance api is "up"
    And sns service is "up"
    And I receive message with summary equals "STRONG_BUY"
    And there are 6 clients available in DB
    When handler is triggered
    Then there should be 6 messages sent via sns
    And sns messages payload should have all client_id's got from clients table
    And sns messages payload symbol should be equal "STRONG_BUY"
    And sns messages payload operation should be equal "BUY"
    And process should exit with 0

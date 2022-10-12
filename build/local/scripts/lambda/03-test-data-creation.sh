#!/bin/bash


echo "-----------------Script-03----------------- [validator]"

echo "########### Inserting test client on DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb put-item \
    --endpoint-url=http://localstack:4566 \
    --table-name crypto_robot.clients \
    --item '{
              "client_id": {
                "S": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"
              },
              "active": {
                "BOOL": true
              },
              "locked_until": {
                "S": "2022-09-17T12:05:07.45066-03:00"
              },
              "locked": {
                "BOOL": false
              },
              "cash_available": {
                "N": "10000.00"
              },
              "cash_amount": {
                "N": "1000.00"
              },
              "cash_reserved": {
                "N": "0"
              },
              "crypto_available": {
                "N": "1.000000"
              },
              "crypto_amount": {
                "N": "1.000000"
              },
              "crypto_symbol": {
                "S": "BTC"
              },
              "crypto_reserved": {
                "N": "0"
              },
              "buy_on": {
                "N": "1"
              },
              "sell_on": {
                "N": "1"
              },
              "ops_timeout_seconds": {
                "N": "60"
              },
              "operation_stop_loss": {
                "N": "50"
              },
              "day_stop_loss": {
                "N": "500"
              },
              "month_stop_loss": {
                "N": "500"
              },
              "operation_amount_percentage": {
                "N": "5.5"
              },
              "month_sell_cap": {
                "N": "25000"
              },
              "symbols": {
                "L": [
                  {
                    "S": "BTC"
                  },
                  {
                    "S": "SOL"
                  }
                ]
              },
              "monthly_summary": {
                "M": {
                  "month": {
                    "S": "08/2022"
                  },
                  "amount_sold": {
                    "N": "23000.42"
                  },
                  "amount_bought": {
                    "N": "37123.42"
                  },
                  "profit": {
                    "N": "1032.32"
                  },
                  "crypto": {
                    "L": [
                      {
                        "M": {
                          "profit": {
                            "N": "-53"
                          }
                        }
                      }
                    ]
                  }
                }
              },
              "daily_summary": {
                "M": {
                  "day": {
                    "S": "14/08/2022"
                  },
                  "amount_sold": {
                    "N": "23000.42"
                  },
                  "amount_bought": {
                    "N": "37123.42"
                  },
                  "profit": {
                    "N": "-53"
                  },
                  "crypto": {
                    "L": [
                      {
                        "M": {
                          "profit": {
                            "N": "-53"
                          }
                        }
                      }
                    ]
                  }
                }
              }
            }' \
    --return-consumed-capacity TOTAL

echo "########### Inserting test client credentials on DynamoDB 'crypto_robot.credentials' table ###########"
aws dynamodb put-item \
    --endpoint-url=http://localstack:4566 \
    --table-name crypto_robot.credentials \
    --item '{
              "client_id": {
                "S": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"
              },
              "api_key": {
                "S": "febf7e0f-7ef2-4da8-9afe-7f9f3721d9c6"
              },
              "api_secret": {
                "S": "a7aca6d4f67519fbb4dc65b159b4e9526b069a2cb5f515d4690bce05ba81e6e5967f477e0ce3affa7c80843f3efed1cee9b0c062"
              }
            }' \
    --return-consumed-capacity TOTAL

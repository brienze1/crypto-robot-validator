#!/bin/sh

echo "########### Inserting test client on DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb put-item \
    --endpoint-url=http://localhost:4566 \
    --profile localstack \
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
              "cash_amount": {
                "N": "1000.00"
              },
              "cash_reserved": {
                "N": "0"
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

#echo "########### Get created item from table ###########"
#aws dynamodb get-item \
#    --endpoint-url=http://localhost:4566 \
#    --profile localstack \
#    --table-name crypto_robot.clients \
#    --key '{ "client_id": { "S": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e" } }'
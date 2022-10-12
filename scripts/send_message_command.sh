#!/bin/sh

echo "########### Sending message to SNS ###########"
aws sns publish \
--endpoint-url=http://localhost:4566 \
--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoOperationTriggerTopic \
--profile localstack \
--message '{
             "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
             "operation": "BUY",
             "symbol": "BTC",
             "start_time": "2022-09-17T12:05:07.45066-03:00"
           }'

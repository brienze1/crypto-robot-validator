#!/bin/sh

echo "########### Sending message to SNS ###########"
aws sns publish \
--endpoint-url=http://localhost:4566 \
--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoAnalysisSummaryTopic \
--profile localstack \
--message '{
             "summary": "STRONG_BUY",
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
           }'

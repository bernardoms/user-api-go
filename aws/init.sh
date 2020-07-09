#!/bin/bash
set -x
echo "runing"
awslocal sns create-topic --name user_update_notify
awslocal sqs create-queue --endpoint-url=http://localstack:4576 --queue-name user_notify_queue_1;
awslocal --endpoint-url=http://localstack:4575 sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:user_update_notify --protocol sqs --notification-endpoint http://localstack:4576/queue/user_notify_queue_1
set +x
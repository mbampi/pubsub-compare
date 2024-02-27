#!/bin/bash

if [ "$MESSAGE_BROKER" == "kafka" ]; then
  docker-compose -f docker-compose.yml -f docker-compose.kafka.yml up -d
elif [ "$MESSAGE_BROKER" == "rabbitmq" ]; then
  docker-compose -f docker-compose.yml -f docker-compose.rabbitmq.yml up -d
else
  echo "Please set the MESSAGE_BROKER environment variable to either 'kafka' or 'rabbitmq'."
fi

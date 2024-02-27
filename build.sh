# Description: Build the docker containers for the application and the message broker.

if [ "$MESSAGE_BROKER" == "kafka" ]; then
  docker-compose -f docker-compose.yml -f docker-compose.rabbitmq.yml down
  docker-compose -f docker-compose.yml -f docker-compose.kafka.yml down
  docker-compose -f docker-compose.yml -f docker-compose.kafka.yml build
elif [ "$MESSAGE_BROKER" == "rabbitmq" ]; then
  docker-compose -f docker-compose.yml -f docker-compose.rabbitmq.yml down
  docker-compose -f docker-compose.yml -f docker-compose.kafka.yml down
  docker-compose -f docker-compose.yml -f docker-compose.rabbitmq.yml build 
else
  echo "Please set the MESSAGE_BROKER environment variable to either 'kafka' or 'rabbitmq'."
fi

# Event Streaming Comparison

This repository contains the code and data for the event streaming comparison. The comparison is based on the following event streaming systems:
- Apache Kafka
- RabbitMQ

The comparison is based on the following criteria:
- Throughput
- Latency
- Resource consumption (CPU, memory)
- Error rate

Message used for the comparison:
- Around 120-150 bytes
- JSON format
- Similar to an IoT message

```json
{
  "id": 1,
  "timestamp": "2020-01-01T00:00:00.000Z",
  "temperature": 20.0,
  "humidity": 50.0,
  "pressure": 1013.25
}
```

## How to Run

go to `docker-compose.yml` and change the `BROKER` and the `BROKER_ADDRESS` environment variable and then run the following commands:
 
```bash
export MESSAGE_BROKER=kafka # or rabbitmq

./build.sh && ./run.sh
 ```

### What is it

Docker-compose for Kafka + Zookeeper. Based on [wurstmeister/kafka-docker](https://github.com/wurstmeister/kafka-docker).

### Running

Starting the Stack:

`docker-compose up -d`

Creating a topic:

`./kafka_2.11-2.0.0/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test`

Send some messages:

`./kafka_2.11-2.0.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test`

Start a consumer:

`./kafka_2.11-2.0.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning`
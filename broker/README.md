### What is it

Docker-compose for a Message Broker service implemented using Kafka + Zookeeper. 
Based on [wurstmeister/kafka-docker](https://github.com/wurstmeister/kafka-docker).

### Running

Starting the Stack:

`docker-compose up -d`

Creating a topic:

`./kafka_2.11-2.0.0/bin/kafka-topics.sh --zookeeper localhost:2181 --create --topic test --replication-factor 1 --partitions 1`

List available topics:

`./kafka_2.11-2.0.0/bin/kafka-topics.sh --zookeeper localhost:2181 --list`

Delete a topic:

`./kafka_2.11-2.0.0/bin/kafka-topics.sh --zookeeper localhost:2181 --delete --topic test`

Send some messages:

`./kafka_2.11-2.0.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test`

Start a consumer:

`./kafka_2.11-2.0.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning`
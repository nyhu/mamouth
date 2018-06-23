# MAMOUTH real time data archiving

## Introduction


## KAFKA TOPICS MANAGEMENT

Kafka without docker force config file usage, so for cluster deployment we'll have to template kafka config files and have a specific process to handle it.
Compressed kafka source of the same version was added to the repository to keep version tracking if we wanted to come back to a traditionnal deployment.
Confluent is the main company building arround kafka and add several layers that can be proven useful. They are the only one maintaining a kafka docker container on a professionnal maner. Confluent and kafka version tracking doesn't match. 

For testing purposes, feel free to install `kafkacat` with your own package manager.
With it you'll be able to produce, consume and list topics. Several option are available, note that command flag --help is more complete than manual.
Local installation without docker is the esyest way to understand kafka. Installation process is describe bellow.

### Creation

You can then use `create-topic.sh` to create desired topics.

Topic creation is the first step 
`bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic my-metron-topic`

### List

`bin/kafka-topics.sh --list --zookeeper localhost:2181`

## Produce

`bin/kafka-console-producer.sh --broker-list localhost:9092 --topic my-metron-topic`
Then write as many message as you want, or let the readline open for further testing

## Consume

`bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-metron-topic`
By default it will consume from last offset before ending. If you wanna read from offset 0, add `--from-beginning`


## Dependencies for cluster

- Docker
- Docker-machine
- Docker-compose

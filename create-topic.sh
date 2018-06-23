#!/bin/bash

if [ -z "$1" ]
  then
    echo "No argument supplied"
    exit 1
fi

ZK_LEADER_URL=51.15.231.63:2181

docker run \
    --net=host \
    --rm \
    confluentinc/cp-kafka:4.1.0 \
    kafka-topics --create --topic $1 --partitions 9 --replication-factor 1 --if-not-exists --zookeeper $ZK_LEADER_URL

docker run \
    --net=host \
    --rm \
    confluentinc/cp-kafka:4.1.0 \
    kafka-topics --describe --topic $1 --zookeeper $ZK_LEADER_URL

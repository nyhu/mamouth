import json, time, random
from pykafka import KafkaClient

client = KafkaClient(hosts="51.15.231.63:29092")
print "actual kafka topics: ", client.topics

topic = client.topics['test']

consumer = topic.get_simple_consumer()

for message in consumer:
    if message is not None:
        print message.offset, message.value

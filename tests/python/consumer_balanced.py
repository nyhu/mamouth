import json, time, random
from pykafka import KafkaClient

client = KafkaClient(hosts="51.15.231.63:29092")
print "actual kafka topics: ", client.topics

#topic = client.topics['metron-test']
topic = client.topics['test']

consumer = topic.get_balanced_consumer(
    consumer_group='testgroup',
    auto_commit_enable=True,
    zookeeper_connect='163.172.102.78:2181'
)

for message in consumer:
    if message is not None:
        #print message.offset, message.value
        print message.value

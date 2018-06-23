import json, time, random, Queue
from pykafka import KafkaClient

client = KafkaClient(hosts="51.15.231.63:29092")
print "actual kafka topics: ", client.topics

topic = client.topics['test']

def makeMessage():
    return json.dumps({
            "sensor_name":"[test]simulator/mamouth",
            "time": int(time.time() * 1000),
            "value": str(random.uniform(-3000, +3000)),
        })

print "exemple test message: ", makeMessage()

#synch:
# with topic.get_sync_producer() as producer:
#     for i in range(9):
#         producer.produce(makeOldMessage())

#asynch:
nb_message = 0
with topic.get_producer(delivery_reports=True) as producer:
    while nb_message < 500000:
        for partion_nb in range(0, 9):
            producer.produce(makeMessage(), partition_key='{}'.format(partion_nb))
            nb_message += 1
        while True:
            try:
                msg, exc = producer.get_delivery_report(block=False)
                if exc is not None:
                    print 'Failed to deliver msg {}: number {} on partition {}'.format(
                        nb_message, msg.partition_key, repr(exc))
                else:
                    print 'Successfully delivered msg number {} on partition {}'.format(
                    nb_message, msg.partition_key)
            except Queue.Empty:
                break

print "msg sent = " + str(nb_message)

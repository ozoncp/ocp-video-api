# test = curl -X POST -k http://localhost:7000/v1/videos -d '{ "slide_id": 555, "link": "/link/555"}'

from kafka import KafkaConsumer
import json

consumer = KafkaConsumer(
    'video',
     bootstrap_servers=['localhost:9094'],
     auto_offset_reset='earliest',
     enable_auto_commit=True,
     group_id='my-group',
     value_deserializer=lambda x: json.loads(x.decode('utf-8')))

for msg in consumer:
    print (msg)

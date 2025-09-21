"""
Backend module for processing telemetry data.
"""

from os import getenv
from sys import exit as sys_exit

from kafka.errors import NoBrokersAvailable

from utils.constants import KAFKA_DEFAULT_ADDRESS, KAFKA_DEFAULT_PORT
from utils.kafka_consumer import TelemetryConsumer


if __name__ == "__main__":
    args = {
        'kafka_address': getenv('KAFKA_ADDRESS', KAFKA_DEFAULT_ADDRESS),
        'kafka_port': getenv('KAFKA_PORT', KAFKA_DEFAULT_PORT),
    }

    print('Starting consumer on', args['kafka_address'], args['kafka_port'])

    try:
        consumer_obj = TelemetryConsumer(args['kafka_address'], args['kafka_port'])
    except NoBrokersAvailable:
        print("No Kafka brokers available!")
        sys_exit(1)

    try:
        consumer_obj.subscribe_to_pattern(r'^telemetry\..+')
        msg_consumer = consumer_obj.get_consumer()
        # Receive loop
        for message in msg_consumer:
            print("Received message on topic", message.topic)
    finally:
        consumer_obj.close()

"""
Kafka consumer for processing telemetry data.
"""
from collections.abc import Set
import json
from kafka import KafkaConsumer

from utils.constants import KAFKA_DEFAULT_ADDRESS, KAFKA_DEFAULT_PORT


class TelemetryConsumer:
    """
    Kafka consumer for processing telemetry data.

    :param address: Kafka broker address.
    :param port: Kafka broker port.
    """
    def __init__(self, address=KAFKA_DEFAULT_ADDRESS, port=KAFKA_DEFAULT_PORT):
        self.consumer = KafkaConsumer(
            bootstrap_servers=f"{address}:{port}",
            auto_offset_reset='earliest',
            enable_auto_commit=True,
            group_id='python-reader',
            value_deserializer=lambda x: json.loads(x.decode('utf-8')),
            allow_auto_create_topics=False,
        )

    def subscribe_to_pattern(self, pattern: str):
        """
        Subscribe to Kafka topics via a regex pattern.

        :param pattern: Regex pattern to subscribe to topics.
        :raises ValueError: If pattern is empty.
        """
        if not pattern:
            raise ValueError("Pattern must not be empty.")
        self.consumer.subscribe(pattern=pattern)

    def subscribe_to_topics(self, topics: Set[str]):
        """
        Subscribe to Kafka topics.

        :param topics: Set of topics to subscribe to.
        :raises ValueError: If topics set is empty.
        """
        if not topics:
            raise ValueError("Topics must not be empty.")
        self.consumer.subscribe(topics=topics)

    def get_consumer(self):
        """Get the Kafka consumer instance."""
        return self.consumer

    def close(self):
        """Close the Kafka consumer."""
        self.consumer.close()

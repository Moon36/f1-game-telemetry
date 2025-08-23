"""
Kafka consumer for processing telemetry data.
"""

import json
from kafka import KafkaConsumer

from utils.constants import KAFKA_DEFAULT_ADDRESS, KAFKA_DEFAULT_PORT


class TelemetryConsumer:
    """Kafka consumer for processing telemetry data."""
    def __init__(self, address=KAFKA_DEFAULT_ADDRESS, port=KAFKA_DEFAULT_PORT):
        self.consumer = KafkaConsumer(
            bootstrap_servers=f"{address}:{port}",
            auto_offset_reset='earliest',
            enable_auto_commit=True,
            group_id='python-reader',
            value_deserializer=lambda x: json.loads(x.decode('utf-8')),
            allow_auto_create_topics=False,
        )

    def subscribe(self, topics: list|None=None, pattern: str=''):
        """Subscribe to Kafka topics or patterns."""
        if topics is None:
            topics = []
        if (not topics and not pattern) or (topics and pattern):
            raise ValueError("Must provide either topics or pattern to subscribe.")
        self.consumer.subscribe(topics=topics, pattern=pattern)

    def get_consumer(self):
        """Get the Kafka consumer instance."""
        return self.consumer

    def close(self):
        """Close the Kafka consumer."""
        self.consumer.close()

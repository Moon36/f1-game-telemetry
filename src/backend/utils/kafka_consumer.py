"""
Kafka consumer for processing telemetry data.
"""
from collections.abc import Sequence
import json
from kafka import KafkaConsumer

from utils.constants import KAFKA_DEFAULT_ADDRESS, KAFKA_DEFAULT_PORT, KAFKA_DEFAULT_GROUP_ID


class TelemetryConsumer:
    """
    Kafka consumer for processing telemetry data.

    :param address: Kafka broker address.
    :param port: Kafka broker port.
    :param group_id: Kafka consumer group ID.
    """
    def __init__(self,
                 address: str=KAFKA_DEFAULT_ADDRESS,
                 port: str | int=KAFKA_DEFAULT_PORT,
                 group_id: str=KAFKA_DEFAULT_GROUP_ID):
        self.consumer = KafkaConsumer(
            bootstrap_servers=f"{address}:{port}",
            auto_offset_reset='earliest',
            enable_auto_commit=True,
            group_id=group_id,
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

    def subscribe_to_topics(self, topics: Sequence[str]):
        """
        Subscribe to Kafka topics.

        :param topics: Sequence of topics to subscribe to.
        :raises ValueError: If topics sequence is empty.
        """
        if not any(topics):
            raise ValueError("Topics must not be empty.")
        self.consumer.subscribe(topics=topics)

    def get_consumer(self) -> KafkaConsumer:
        """Get the Kafka consumer instance."""
        return self.consumer

    def close(self):
        """Close the Kafka consumer."""
        self.consumer.close()

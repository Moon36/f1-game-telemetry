"""Message translator for converting topics into specific message types."""
from typing import Type, Dict
from messages.llm_messages.base_message import LLMBaseMessage
from messages.llm_messages.event_message import EventMessage
from database.database import SimpleMemoryDatabase


class MessageTranslator:
    """Translates message topics into specific LLMBaseMessage instances."""

    __message_registry__: Dict[str, Type[LLMBaseMessage]] = {
        "telemetry.event": EventMessage,
    }

    @classmethod
    def translate(cls, topic: str, *args, **kwargs) -> LLMBaseMessage:
        """
        Translate a topic and arguments into a specific message instance.
        
        :param topic: The message topic.
        :param args: Positional arguments for the message constructor.
        :param kwargs: Keyword arguments for the message constructor.
        :return: An instance of LLMBaseMessage corresponding to the topic.
        :raises ValueError: If the topic is not registered.
        """
        topic_lower = topic.lower()
        if topic_lower not in cls.__message_registry__:
            raise ValueError(f"Unknown message topic: {topic}")

        message_class = cls.__message_registry__[topic_lower]
        # Beware of message header. Header is (probably) irrelevant for the LLM message.
        return message_class(*args, **kwargs)

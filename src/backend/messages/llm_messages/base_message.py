"""Basic message class for LLM messages."""
from abc import ABC
from dataclasses import dataclass

from messages.llm_messages.message_priorities import MessagePriority

@dataclass
class LLMBaseMessage(ABC):
    """Base class for LLM messages."""
    message_priority: MessagePriority
    message_type: str
    message_content: str
    data: dict


    @property
    def priority(self) -> MessagePriority:
        """Get the message priority."""
        return self.message_priority


    @property
    def type(self) -> str:
        """Get the message type."""
        return self.message_type


    @property
    def message(self) -> str:
        """Get the message content to be sent to the LLM."""
        return f'{self.message_content}\nData: {self.data}'


    def __str__(self) -> str:
        return self.message

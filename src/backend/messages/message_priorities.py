"""Defines message priority levels for processing in the priority queue."""
from enum import IntEnum

class MessagePriority(IntEnum):
    """Message priority levels for processing in the priority queue."""
    HIGH = 1
    MEDIUM = 5
    LOW = 10

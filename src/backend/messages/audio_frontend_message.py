"""Provides the layout for the message to exchange audio data between the backend and frontend."""

from base64 import b64encode
from dataclasses import dataclass
from json import dumps
from collections.abc import Sequence

from piper.voice import AudioChunk

@dataclass
class AudioChunkMessage:
    """The AudioChunkMessage class represents an audio message exchanged between the backend and frontend."""

    audio_chunks: Sequence[AudioChunk]

    def to_dict(self) -> dict:
        """Convert the AudioMessage to a dictionary."""
        return {
            'sample_rate': list(map(lambda chunk: chunk.sample_rate, self.audio_chunks)),
            'channels': list(map(lambda chunk: chunk.sample_channels, self.audio_chunks)),
            'audio_data': list(map(lambda chunk: b64encode(chunk.audio_int16_bytes).decode('utf-8'),
                                         self.audio_chunks))
        }

    def to_json(self) -> str:
        """Convert the AudioMessage to a JSON-serializable dictionary."""
        return dumps(self.to_dict())

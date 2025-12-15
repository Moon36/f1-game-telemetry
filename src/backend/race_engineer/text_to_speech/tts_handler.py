"""This file handles the interactions with the text-to-speech model."""

from collections.abc import Sequence

from piper import PiperVoice
from piper.voice import AudioChunk


class TTSHandler:
    """Handles text-to-speech synthesis using Piper."""


    def __init__(self, voice_path: str):
        """
        Initialize the TTS handler with a specific Piper voice.
        
        Args:
            voice: An instance of PiperVoice representing the desired voice.
        """
        self.voice_path = voice_path
        self.voice = self.__load_voice__(voice_path)


    def __load_voice__(self, voice_name: str) -> PiperVoice:
        return PiperVoice.load(voice_name)

    def synthesize_text(self, text: str) -> Sequence[AudioChunk]:
        """
        Synthesize speech from text and return audio chunks.
        
        :param text: The text to be synthesized.
        :return: An iterable of AudioChunk objects representing the synthesized speech.
        """
        return list(self.voice.synthesize(text))

"""This file provides the race engineer class."""

from collections.abc import Sequence
from threading import Lock

from piper.voice import AudioChunk
from requests import HTTPError

from race_engineer.text_generation.config.llm_sys_prompts import Personalities
from race_engineer.text_generation.llm_handler import LLMHandler
from race_engineer.text_to_speech.tts_handler import TTSHandler
from messages.llm_messages.base_message import LLMBaseMessage

class RaceEngineer:
    """The Race Engineer singleton class combines LLM and TTS functionalities."""

    _instance = None
    _initialized = False
    _lock = Lock()

    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            with cls._lock:
                cls._instance = super(RaceEngineer, cls).__new__(cls)
                cls._instance._initialized = False
        return cls._instance


    def __init__(self,
                 model: str,
                 base_url: str,
                 voice_path: str,
                 system_instruction: str | None = None,
                 personality: Personalities = Personalities.NEUTRAL):
        if self._initialized:
            # TODO: log warning
            return

        self.model = model
        self.base_url = base_url
        self.voice_path = voice_path
        self.system_instruction = system_instruction
        self.personality = personality

        self.llm_handler = LLMHandler(
            base_url=self.base_url,
            model_name=self.model,
            system_instruction=self.system_instruction,
            personality=self.personality
        )
        self.tts_handler = TTSHandler(voice_path=self.voice_path)

        self._initialized = True


    def generate_radio_message(self, message: LLMBaseMessage | str) -> Sequence[AudioChunk]:
        """
        Generates a radio message by combining LLM text generation and TTS synthesis.

        :param message: The message to generate text and speech for.
        :return: The generated speech audio in bytes.
        """
        try:
            generated_text = self.llm_handler.send_message(str(message))
        except (TimeoutError, HTTPError, ValueError) as e:
            print('LLM handler returned an error:', e)
            return []
        print('RE:', 'Generated text from LLM:\n', generated_text)
        speech_audio = self.tts_handler.synthesize_text(generated_text)
        return speech_audio

"""This file handles interactions with the LLM for text generation."""

import requests
from urllib.parse import urljoin

from race_engineer.text_generation.config.llm_sys_prompts import Personalities
from race_engineer.text_generation.utils import constants


class LLMHandler:
    """Handles communication with the LLM for generating text messages."""

    def __init__(self, base_url: str,
                 model_name: str,
                 system_instruction: str | None = None,
                 personality: Personalities = Personalities.NEUTRAL):
        self.base_url = base_url
        self.model_name = model_name
        self.system_instruction = system_instruction
        self.personality = personality
        self.url_chat_endpoint = urljoin(self.base_url, constants.LLM_CHAT_ENDPOINT)

        self.conversation = self.__prepare_conversation__()

        print('LLM Configuration:', self.__get_config__())

    def __prepare_conversation__(self) -> dict:
        """Prepares the conversation payload for the LLM request."""
        return {
            'model': self.model_name,
            'messages': self.__set_system_message__()
        }

    def __set_system_message__(self) -> list[dict[str, str]]:
        """Prepares the system instruction by combining role and personality.

        :return: The combined system instruction string.
        """
        content = self.system_instruction if self.system_instruction is not None else ''
        content += '\n' + self.personality.value

        return [{
            "role": "system",
            "content": content
        }]

    def __get_config__(self) -> dict:
        """Returns the current LLM handler configuration."""
        return {
            'base_url': self.base_url,
            'chat_endpoint': self.url_chat_endpoint,
            'model_name': self.model_name,
            'system_instruction': self.system_instruction,
            'personality': self.personality
        }

    def send_message(self, message: str, include_history: bool = True) -> str:
        """
        Generates a message using the LLM based on the user prompt.
        
        :param message: The user prompt to send to the LLM.
        :return: The generated message from the LLM.
        :raises TimeoutError: If the request to the LLM times out.
        :raises HTTPError: If the LLM returns a non-200 status code.
        :raises ValueError: If the LLM response cannot be parsed.
        """
        user_message = {
            "role": "user",
            "content": message
        }
        self.conversation['messages'].append(user_message)

        conv = self.conversation
        if not include_history:
            conv = self.__prepare_conversation__()
            conv['messages'].append(user_message)

        try:
            response = requests.post(self.url_chat_endpoint, json=conv, timeout=constants.REQUEST_TIMEOUT)
        except requests.exceptions.ConnectionError as e:
            raise ConnectionError(f"Could not connect to LLM service. Is it running on {self.base_url}?") from e
        except requests.exceptions.Timeout as e:
            raise TimeoutError("LLM request timed out") from e

        response.raise_for_status()

        reply = self.__parse_response__(response)
        self.conversation['messages'].append({
            "role": "assistant",
            "content": reply
        })

        return reply

    def __parse_response__(self, response: requests.Response) -> str:
        """Parses the LLM response to extract the generated message.

        :param response: The response object from the LLM.
        :return: The generated message from the LLM.
        :raises ValueError: If the response cannot be parsed as JSON.
        """
        try:
            data = response.json()
        except requests.exceptions.JSONDecodeError as e:
            raise ValueError("Failed to decode LLM response as JSON! Is the LLM endpoint broken?") from e

        if 'choices' not in data or not data['choices']:
            raise ValueError("LLM response JSON has unexpected format: 'choices' key missing or empty")
        if 'message' not in data['choices'][0] or 'content' not in data['choices'][0]['message']:
            raise ValueError("LLM response JSON has unexpected format: 'message' or 'content' key missing")

        return data['choices'][0]['message']['content']

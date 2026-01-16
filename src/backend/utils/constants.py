"""
Provides constants for the backend module.
"""
KAFKA_DEFAULT_ADDRESS = 'localhost'
KAFKA_DEFAULT_PORT = '9092'
KAFKA_TOPIC_PATTERN = r'^telemetry\..+'
KAFKA_DEFAULT_GROUP_ID = 'backend-python-consumer-group'

BACKEND_DEFAULT_PORT = '8282'
RACE_ENGINEER_DEFAULT_PORT = '8283'

LLM_DEFAULT_BASE_URL = 'http://model-runner.docker.internal/engines/v1/'
LLM_DEFAULT_MODEL = 'ai/smollm2'
RACE_ENGINEER_DEFAULT_VOICE = './race_engineer/text_to_speech/voices/en_US-arctic-medium.onnx'

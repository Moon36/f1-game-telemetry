"""
Backend module for processing telemetry data.

TODO: Add race engineer conditional if RE should be started or not.
"""
import asyncio
import json
from os import getenv
from sys import exit as sys_exit
from collections.abc import Sequence

from kafka.errors import NoBrokersAvailable

from database.database import SimpleMemoryDatabase
from messages.participants_packet import ParticipantsPacket
from race_engineer.race_engineer import RaceEngineer
from race_engineer.text_generation.config import llm_sys_prompts
from race_engineer.message_scheduler import LLMMessageScheduler
from utils import constants
from utils.kafka_consumer import TelemetryConsumer
from utils.ws_server import WebSocketServer
from utils.message_translator import MessageTranslator
from utils.packet_translator import PacketTranslator


def __message_consumer_forward_task__(
    consumer_obj: TelemetryConsumer,
    kafka_topic_pattern: str,
    loop: asyncio.AbstractEventLoop,
    ws_server: WebSocketServer):
    """
    This method is a blocking method that consumes messages from Kafka. The messages are then broadcasted to all
    connected WebSocket clients using the provided WebSocketServer instance.
    
    :param consumer_obj: An instance of TelemetryConsumer to consume messages from Kafka.
    :param kafka_topic_pattern: The Kafka topic pattern to subscribe to.
    :param loop: The main event loop to schedule the async broadcast tasks.
    :param ws_server: An instance of WebSocketServer to broadcast messages to clients.
    """
    try:
        consumer_obj.subscribe_to_pattern(kafka_topic_pattern)
        msg_consumer = consumer_obj.get_consumer()
        print(f'Consumer subscribed to topics matching pattern: {kafka_topic_pattern}')

        for record in msg_consumer:
            payload = json.dumps({'topic': record.topic, 'data': record.value})
            loop.call_soon_threadsafe(lambda p=payload: loop.create_task(ws_server.broadcast_message(p)))
    finally:
        consumer_obj.close()


def __message_consumer_re_task__(
    consumer_obj: TelemetryConsumer,
    kafka_topic_list: Sequence[str],
    scheduler: LLMMessageScheduler):
    """
    This method is a blocking method that consumes messages from Kafka. The messages are then processed by the AI
    race engineer and broadcasted to all connected WebSocket clients using the provided WebSocketServer instance.
    
    :param consumer_obj: An instance of TelemetryConsumer to consume messages from Kafka.
    :param kafka_topic_pattern: The Kafka topic pattern to subscribe to.
    :param loop: The main event loop to schedule the async broadcast tasks.
    :param ws_server: An instance of WebSocketServer to broadcast messages to clients.
    """
    try:
        consumer_obj.subscribe_to_topics(kafka_topic_list)
        msg_consumer = consumer_obj.get_consumer()
        print(f'RE Consumer subscribed to topics matching pattern: {kafka_topic_list}')

        for record in msg_consumer:
            #print("Race Engineer task received message on topic", record.topic, 'with value', record.value)
            try:
                message = MessageTranslator.translate(record.topic, **record.value)
            except ValueError as e:
                print('Message translation error:', e)
                continue

            if message.message_type == 'EVENT':
                if 'event_code' in message.data:
                    if message.data['event_code'] == 'FTLP'\
                        or message.data['event_code'] == 'RTMT'\
                        or message.data['event_code'] == 'DRSE'\
                        or message.data['event_code'] == 'DRSD'\
                        or message.data['event_code'] == 'TMPT'\
                        or message.data['event_code'] == 'RCWN'\
                        or message.data['event_code'] == 'PENA'\
                        or message.data['event_code'] == 'LGOT':
                        #or message.data['event_code'] == 'OVTK':
                        print('RELEVANT event_code:', message.data['event_code'])

                        # Run heavy LLM and TTS processing in a separate thread
                        scheduler.schedule_message(message.priority, message)
    except Exception as e:
        print('Some error occurred in RE consumer task:', e)
    finally:
        consumer_obj.close()
        scheduler.stop()

    # DELME: Temporary loop to simulate RE messages for testing
    #from time import sleep
    #from messages.audio_frontend_message import AudioChunkMessage
    #while True:
    #    sleep(10)
#
    #    message = MessageTranslator.translate('telemetry.event',
    #                                          **{
    #                                              'M_eventStringCode': [66, 85, 84, 78],
    #                                              'M_eventDetails': {'ButtonStatus': 16}
    #                                            })
    #    print('Sending message', f'\"{message}\"', 'to Race Engineer')
    #    tts_res = race_engineer.generate_radio_message(message)
    #    if any(tts_res):
    #        audio_message = AudioChunkMessage(tts_res)
    #        loop.call_soon_threadsafe(
    #            lambda p=audio_message: loop.create_task(ws_server.broadcast_message(p.to_json()))
    #        )
    #    else:
    #        print('No TTS response from Race Engineer')


async def main():
    """Main async function to coordinate WebSocket server and Kafka consumer."""
    args = {
        'kafka_address': getenv('KAFKA_ADDRESS', constants.KAFKA_DEFAULT_ADDRESS),
        'kafka_port': getenv('KAFKA_PORT', constants.KAFKA_DEFAULT_PORT),
        'ws_port': getenv('BACKEND_PORT', constants.BACKEND_DEFAULT_PORT),
        're_ws_port': getenv('RE_BACKEND_PORT', constants.RACE_ENGINEER_DEFAULT_PORT),
        'llm_base_url': getenv('LLM_BASE_URL', constants.LLM_DEFAULT_BASE_URL),
        'llm_model': getenv('LLM_MODEL', constants.LLM_DEFAULT_MODEL),
        'race_engineer_voice': getenv('RACE_ENGINEER_VOICE', constants.RACE_ENGINEER_DEFAULT_VOICE)
    }

    print('Starting consumer on', args['kafka_address'], args['kafka_port'])

    try:
        consumer_obj = TelemetryConsumer(args['kafka_address'], args['kafka_port'])
        re_consumer_obj = TelemetryConsumer(args['kafka_address'], args['kafka_port'], 'race-engineer-consumer-group')
    except NoBrokersAvailable:
        print("No Kafka brokers available!")
        sys_exit(1)

    print(f'Starting WebSocket server on ws://0.0.0.0:{args["ws_port"]}')
    ws_server = WebSocketServer(host='0.0.0.0', port=int(args['ws_port']))
    ws_server_task = asyncio.create_task(ws_server.start())
    print(f'Starting WebSocket server for AI Race Engineer on ws://0.0.0.0:{args["re_ws_port"]}')
    re_ws_server = WebSocketServer(host='0.0.0.0', port=int(args['re_ws_port']))
    re_ws_server_task = asyncio.create_task(re_ws_server.start())

    race_engineer = RaceEngineer(args['llm_model'],
                                 args['llm_base_url'],
                                 constants.RACE_ENGINEER_DEFAULT_VOICE,
                                 llm_sys_prompts.LLM_SYS_PROMPT,
                                 llm_sys_prompts.Personalities.ANGRY)

    event_loop = asyncio.get_event_loop()

    llm_msg_scheduler = LLMMessageScheduler(event_loop, re_ws_server, race_engineer, cooldown_time=5)
    llm_msg_scheduler.start()

    # Make consumer forward thread
    consumer_forward_thread = asyncio.to_thread(__message_consumer_forward_task__,
                                                consumer_obj,
                                                constants.KAFKA_TOPIC_PATTERN,
                                                event_loop,
                                                ws_server)

    race_engineer_thread = asyncio.to_thread(__message_consumer_re_task__,
                                             re_consumer_obj,
                                             ['telemetry.event'],
                                             llm_msg_scheduler)

    await asyncio.gather(
        ws_server_task,
        re_ws_server_task,
        consumer_forward_thread,
        race_engineer_thread
    )


if __name__ == "__main__":
    asyncio.run(main())

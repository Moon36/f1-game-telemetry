"""
Backend module for processing telemetry data.
"""
import asyncio
import json
from os import getenv
from sys import exit as sys_exit

from kafka.errors import NoBrokersAvailable

from utils import constants
from utils.kafka_consumer import TelemetryConsumer
from utils.ws_server import WebSocketServer


def __message_consumer_blocking_loop__(
    consumer_obj: TelemetryConsumer,
    kafka_topic_pattern: str,
    ws_server: WebSocketServer,
    loop: asyncio.AbstractEventLoop
):
    """
    This method is a blocking method that consumes messages from Kafka. The messages are then
    broadcasted to all connected WebSocket clients using the provided WebSocketServer instance.
    
    :param consumer_obj: An instance of TelemetryConsumer to consume messages from Kafka.
    :param ws_server: An instance of WebSocketServer to broadcast messages to WebSocket clients.
    :param loop: The main event loop to schedule the async broadcast tasks.
    """
    try:
        consumer_obj.subscribe_to_pattern(kafka_topic_pattern)
        msg_consumer = consumer_obj.get_consumer()
        print(f'Consumer subscribed to topics matching pattern: {kafka_topic_pattern}')

        for record in msg_consumer:
            print(f"Received message on topic {record.topic}")
            # Schedule the async broadcast_message coroutine to be run on the main event loop
            loop.call_soon_threadsafe(
                lambda record=record: asyncio.create_task(ws_server.broadcast_message(json.dumps({
                'topic': record.topic,
                'data': record.value
            })))
            )
    finally:
        consumer_obj.close()


async def kafka_consumer_task(consumer_obj: TelemetryConsumer, kafka_topic_pattern: str, ws_server: WebSocketServer):
    """
    Runs the Kafka consumer in a separate thread.
    
    :param consumer_obj: An instance of TelemetryConsumer to consume messages from Kafka.
    :param ws_server: An instance of WebSocketServer to broadcast messages to WebSocket clients
    """
    loop = asyncio.get_event_loop()
    await loop.run_in_executor(
        None,  # Use the default thread pool executor
        __message_consumer_blocking_loop__,
        consumer_obj,
        kafka_topic_pattern,
        ws_server,
        loop
    )


async def main():
    """Main async function to coordinate WebSocket server and Kafka consumer."""
    args = {
        'kafka_address': getenv('KAFKA_ADDRESS', constants.KAFKA_DEFAULT_ADDRESS),
        'kafka_port': getenv('KAFKA_PORT', constants.KAFKA_DEFAULT_PORT),
        'ws_port': getenv('BACKEND_PORT', constants.BACKEND_DEFAULT_PORT)
    }

    print('Starting consumer on', args['kafka_address'], args['kafka_port'])

    try:
        consumer_obj = TelemetryConsumer(args['kafka_address'], args['kafka_port'])
    except NoBrokersAvailable:
        print("No Kafka brokers available!")
        sys_exit(1)

    print(f'Starting WebSocket server on ws://0.0.0.0:{args["ws_port"]}')
    ws_server = WebSocketServer(host='0.0.0.0', port=int(args['ws_port']))

    async with ws_server:
        print("WebSocket server started")
        consumer_task = asyncio.create_task(kafka_consumer_task(consumer_obj, constants.KAFKA_TOPIC_PATTERN, ws_server))
        await consumer_task


if __name__ == "__main__":
    asyncio.run(main())

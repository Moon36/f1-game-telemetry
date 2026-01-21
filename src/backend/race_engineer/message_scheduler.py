"""
Provides a message scheduler singleton for the LLM Race Engineer. The scheduler processes incoming messages via a
priority queue and submits them to the frontend.
"""
import asyncio
import threading
from queue import PriorityQueue
from threading import Lock
from time import time
from dataclasses import dataclass, field
from typing import Any

from utils.ws_server import WebSocketServer
from race_engineer.race_engineer import RaceEngineer
from messages.audio_frontend_message import AudioChunkMessage
from messages.message_priorities import MessagePriority
from messages.llm_messages.base_message import LLMBaseMessage


@dataclass(order=True)
class PrioritizedItem:
    """Dataclass for items in the priority queue."""
    priority: int
    item: Any=field(compare=False)


class LLMMessageScheduler:
    """Singleton class for scheduling LLM messages based on priority."""

    _instance = None
    _initialized = False
    _lock = Lock()


    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            with cls._lock:
                cls._instance = super(LLMMessageScheduler, cls).__new__(cls)
                cls._instance._initialized = False
        return cls._instance


    def __init__(self,
                 event_loop: asyncio.AbstractEventLoop,
                 ws_server: WebSocketServer,
                 race_engineer: RaceEngineer,
                 cooldown_time: float = 5.0):
        if self._initialized:
            # TODO: log warning
            return

        self.event_loop = event_loop
        self.ws_server = ws_server
        self.race_engineer = race_engineer
        self.cooldown_time = cooldown_time

        self._message_queue = PriorityQueue()
        self._last_process_time = 0.0
        self._do_process_queue = True
        self._interrupt_event = threading.Event()

        self._worker_thread = threading.Thread(target=self.__process_queue__, daemon=True)

        self._initialized = True


    def schedule_message(self, priority: MessagePriority, message: LLMBaseMessage):
        """
        Schedule a message with a given priority.

        :param message: The message for the LLM to be scheduled.
        :param priority: The priority of the message.
        """
        if priority == MessagePriority.HIGH:
            self._interrupt_event.set()
        self._message_queue.put(PrioritizedItem(priority, message))


    def start(self):
        """Start the message processing loop. The loop runs in a separate thread."""
        print('Starting LLM Message Scheduler processing thread')
        self._worker_thread.start()


    def stop(self):
        """Stop the message processing loop."""
        print('Stopping LLM Message Scheduler processing thread')
        self._do_process_queue = False
        self._worker_thread.join(timeout=1)


    def __elapsed_since_last_run__(self) -> float:
        return time() - self._last_process_time


    def __process_queue__(self):
        """
        Process messages from the priority queue and submit them to the Race Engineer for generation and TTS synthesis.
        Finally broadcast the generated audio messages to connected WebSocket clients.
        """
        while self._do_process_queue:
            item = self._message_queue.get(block=True, timeout=None)
            priority, message = item.priority, item.item

            if priority != MessagePriority.HIGH:
                elapsed_time = self.__elapsed_since_last_run__()
                if elapsed_time < self.cooldown_time:
                    wait_time = self.cooldown_time - elapsed_time
                    print('Cooling down for', round(wait_time, 2), 'seconds before processing next message.')
                    # Wait for remaining cooldown OR until event is set (on HIGH priority message)
                    interrupted = self._interrupt_event.wait(timeout=wait_time)
                    if interrupted:
                        self._interrupt_event.clear()

            tts_res = self.race_engineer.generate_radio_message(message)
            if not any(tts_res):
                print('No TTS response from Race Engineer')
                continue

            print('Received TTS response')

            audio_message = AudioChunkMessage(tts_res)
            print('Broadcasting audio message to RE clients')
            asyncio.run_coroutine_threadsafe(
                self.ws_server.broadcast_message(audio_message.to_json()),
                self.event_loop
            )

            self._message_queue.task_done()

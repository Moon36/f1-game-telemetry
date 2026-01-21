# disable=too-few-public-methods
"""This file provides the incident message class."""
from collections.abc import Sequence

from messages.llm_messages.base_message import LLMBaseMessage
from messages.message_priorities import MessagePriority

PRIOROTY = MessagePriority.HIGH
TYPE = "EVENT"
MESSAGE_CONTENT = ("You just received information, that an event has occurred on track. You are provided with the "
                   "event type and some details about it. Communicate this information to the driver over the radio.")

DESCRIPTION_MAP = {
    'SSTA': 'Session Started',
    'SEND': 'Session Ended',
    'FTLP': 'A driver has achieved the fastest lap',                # Needs ID translation
    'RTMT': 'A driver has retired',                                 # Needs ID translation
    'DRSE': 'Race control have enabled DRS',
    'DRSD': 'Race control have disabled DRS',
    'TMPT': 'Team mate is in the pits',                             # Needs ID translation
    'CHQF': 'Chequered Flag has been waved',
    'RCWN': 'The race winner was announced',                        # Needs ID translation
    'PENA': 'A penalty has been issued to a driver',                # Needs ID translation
    'SPTP': 'The speed trap has been triggered by fastest speed',   # Needs ID translation
    'STLG': 'The number of start lights currently showing',
    'LGOT': 'Start lights out',
    'DTSV': 'Drive through penalty served by a driver',             # Needs ID translation
    'SGSV': 'Stop go penalty served by a driver',                   # Needs ID translation
    'FLBK': 'Game Flashback was activated',                         # Needs ID translation
    'BUTN': 'A button status on the steering wheel has changed',    # Needs ID translation
    'RDFL': 'Red flag is shown (session stopped)',
    'OVTK': 'An overtake has occurred',                             # Needs ID translation
}

class EventMessage(LLMBaseMessage):
    """A message representing an event that has occurred on track."""

    def __init__(self, **kwargs):
        if 'M_eventStringCode' in kwargs:
            if isinstance(kwargs['M_eventStringCode'], Sequence)\
                and all(isinstance(b, int) for b in kwargs['M_eventStringCode']):
                kwargs['M_eventStringCode'] = self.__translate_event_code__(kwargs['M_eventStringCode'])

            kwargs = self.get_event_description(kwargs)
        super().__init__(PRIOROTY, TYPE, MESSAGE_CONTENT, kwargs)


    def __translate_event_code__(self, code: Sequence[int]) -> str:
        """Translate event code to human-readable description."""
        return ''.join(chr(b) for b in code)


    def get_event_description(self, event_details: dict) -> dict:
        """Enhances the event with detail descriptions."""
        event_code = event_details.get('M_eventStringCode')
        if not event_code or event_code not in DESCRIPTION_MAP:
            return event_details

        event_details['event_description'] = DESCRIPTION_MAP[event_code]

        return event_details

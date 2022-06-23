from enum import Enum, auto


class SupportedContentType(Enum):
    HTTP_REQUEST_CYCLE = auto()
    TERMINAL_RECORDING = auto()
    CODEBLOCK = auto()
    EVENT = auto()
    IMAGE = auto()
    NONE = auto()

    @staticmethod
    def from_str(s: str):
        values: dict[str, SupportedContentType] = {
            "http-request-cycle": SupportedContentType.HTTP_REQUEST_CYCLE,
            "terminal-recording": SupportedContentType.TERMINAL_RECORDING,
            "codeblock": SupportedContentType.CODEBLOCK,
            "event": SupportedContentType.EVENT,
            "image": SupportedContentType.IMAGE,
            "none": SupportedContentType.NONE,
        }
        return values[s]

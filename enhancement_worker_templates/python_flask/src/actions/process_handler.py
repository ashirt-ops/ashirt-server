from request_types import ProcessBody
from constants import SupportedContentType
from .types import ProcessResultDTO


def handle_process(body: ProcessBody) -> ProcessResultDTO:
    """
    handle_process is called when a web request comess in, is validated, and indicates that work
    needs to be done on a piece of evidence
    """
    accepted_types = [
        SupportedContentType.IMAGE,
        SupportedContentType.CODEBLOCK,
        SupportedContentType.EVENT,
        SupportedContentType.HTTP_REQUEST_CYCLE,
        SupportedContentType.TERMINAL_RECORDING,
        SupportedContentType.NONE,
    ]

    if body.content_type in accepted_types:
        return {
            'action': 'processed',
            'content': 'TBD'
        }
    else:
        return {
            'action': 'rejected'
        }

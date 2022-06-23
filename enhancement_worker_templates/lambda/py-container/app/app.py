import os

from constants import SupportedContentType
from message_types import (
    ProcessBody,
    TestBody,
)
from responses import (
    bad_request,
    error_processing,
    process_success,
    reject_evidence,
    test_passed,
)
from services import AShirtRequestsService


def handler(event, context):

    if TestBody.parse_if_valid(event) is not None:
        return test_passed()

    if (body := ProcessBody.parse_if_valid(event)) is not None:
        return handle_process(body)

    return bad_request()


def handle_process(body: ProcessBody):
    # TODO: Handle custom logic here!

    # filter out unprocessable evidence
    if body.content_type != SupportedContentType.IMAGE:
        return reject_evidence()

    try:
        content = do_processing(body)
        return process_success(content)
    except Exception as err:  # blanket except here, as we don't know what to to expect
        return error_processing(str(err))


def do_processing(body: ProcessBody):
    # Create ashirt services instance
    # ashirt_svc = AShirtRequestsService(
    #     os.environ.get('ASHIRT_BACKEND_URL', ''),
    #     os.environ.get('ASHIRT_ACCESS_KEY', ''),
    #     os.environ.get('ASHIRT_SECRET_KEY', '')
    # )

    # TODO: logic goes here too
    return "Everything is working"

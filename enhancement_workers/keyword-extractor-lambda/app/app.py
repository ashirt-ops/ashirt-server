from constants import SupportedContentType
from message_types import (
    EvidenceCreatedBody,
    TestBody,
)
from responses import (
    bad_request,
    error_processing,
    process_success,
    reject_evidence,
    test_passed,
)
from processing import process_content


def __main__(event, context):

    if TestBody.parse_if_valid(event) is not None:
        return test_passed()

    if (body := EvidenceCreatedBody.parse_if_valid(event)) is not None:
        return handle_evidence_created(body)

    return bad_request()


def handle_evidence_created(body: EvidenceCreatedBody):
    # filter out unprocessable evidence
    if body.content_type != SupportedContentType.IMAGE:
        return reject_evidence()

    try:
        content = process_content(body)
        return process_success(content)
    except Exception as err:  # blanket except here, as we don't know what to to expect
        return error_processing(str(err))

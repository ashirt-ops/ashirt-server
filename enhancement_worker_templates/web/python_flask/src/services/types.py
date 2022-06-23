from dataclasses import dataclass
import mimetypes
from typing import Literal, Optional, TypedDict

from constants.supported_content_type import SupportedContentType

HTTP_METHOD = Literal['GET', 'POST', 'PUT', 'DELETE']


class FileData(TypedDict):
    filename: str
    mimetype: str
    content: bytes


@dataclass(frozen=True)
class RequestConfig:
    """
    RequestConfig abstracts a request so that it can be sent via different libraries,
    in case you don't like requests
    """
    method: HTTP_METHOD
    path: str
    body: Optional[bytes | str] = None
    return_type: Literal["json", "raw"] = "json"


# The below are all inputs for various API calls

class CreateOperationInput(TypedDict):
    slug: str
    name: str


class CreateEvidenceInput(TypedDict):
    notes: str
    content_type: SupportedContentType
    tag_ids: list[int]
    file: FileData


class UpdateEvidenceInput(TypedDict):
    notes: Optional[str]
    content_type: SupportedContentType
    add_tag_ids: list[int]
    remove_tag_ids: list[int]
    file: FileData


class UpsertEvidenceMetadata(TypedDict):
    source: str
    body: str
    status: str
    message: Optional[str]
    canProcess: Optional[bool]


class CreateTagInput(TypedDict):
    name: str
    colorName: Optional[str]


def parse_file(filename: str, binary=True):
    method = 'rb' if binary else 'r'
    with open(filename, method) as fh:
        data = fh.read(-1)

    mimetypes.guess_type(filename)

    return FileData(
        filename=filename,
        content=data,
        mimetype="application/octet-stream"
    )

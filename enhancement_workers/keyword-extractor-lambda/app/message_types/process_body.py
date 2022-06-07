from dataclasses import dataclass
from typing import Any, Literal

from constants import SupportedContentType
from helpers import is_literal
from .generic_request_body import GenericRequestBody


@dataclass(repr=False, frozen=True)
class ProcessBody(GenericRequestBody):
    """
    ProcessBody reflects the message received from AShirt when AShirt requests processing
    """
    type: Literal['process']
    evidence_uuid: str
    operation_slug: str
    content_type: SupportedContentType

    def is_valid_instance(self) -> bool:
        return (
            is_literal(self.type, str, 'process')
            and type(self.evidence_uuid) == str
            and type(self.operation_slug) == str
            and type(self.content_type) == SupportedContentType
        )

    @classmethod
    def from_json(cls, data: dict[str, Any]):
        cls.type = data['type']
        cls.evidence_uuid = data['evidenceUuid']
        cls.operation_slug = data['operationSlug']
        cls.content_type = SupportedContentType.from_str(data['contentType'])

        return cls

from dataclasses import dataclass
from typing import Any, Literal

from constants import SupportedContentType
from helpers import is_literal
from .generic_request_body import GenericRequestBody


@dataclass(repr=False, frozen=True)
class EvidenceCreatedBody(GenericRequestBody):
    """
    EvidenceCreatedBody reflects the message received from AShirt when AShirt requests metadata processing
    """
    type: Literal['evidence_created']
    evidence_uuid: str
    operation_slug: str
    content_type: SupportedContentType

    def is_valid_instance(self) -> bool:
        return all([
            is_literal(self.type, str, 'evidence_created'),
            type(self.evidence_uuid) is str,
            type(self.operation_slug) is str,
            type(self.content_type) is SupportedContentType,
        ])

    @classmethod
    def from_json(cls, data: dict[str, Any]):
        cls.type = data['type']
        cls.evidence_uuid = data['evidenceUuid']
        cls.operation_slug = data['operationSlug']
        cls.content_type = SupportedContentType.from_str(data['contentType'])

        return cls

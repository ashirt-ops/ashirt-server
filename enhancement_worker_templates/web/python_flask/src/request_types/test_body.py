from dataclasses import dataclass
from typing import Any, Literal

from helpers import is_literal
from .generic_request_body import GenericRequestBody


@dataclass(repr=False, frozen=True)
class TestBody(GenericRequestBody):
    """
    TestBody reflects the message received from AShirt when AShirt requests testing
    """

    type: Literal['test']

    def is_valid_instance(self) -> bool:
        return is_literal(self.type, str, 'test')

    @classmethod
    def from_json(cls, data: dict[str, Any]):
        cls.type = data['type']

        return cls

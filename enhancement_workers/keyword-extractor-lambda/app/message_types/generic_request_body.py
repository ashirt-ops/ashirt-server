from abc import ABC, abstractmethod
from typing import Any


class GenericRequestBody(ABC):

    @abstractmethod
    def from_json(cls, data: dict[str, Any]):
        pass

    @abstractmethod
    def is_valid_instance(self, data: dict[str, Any]) -> bool:
        pass

    @classmethod
    def parse_if_valid(cls, data: dict[str, Any]):
        """
        parse_if_valid checks that the given data is valid, then parses it.
        if is not valid, or if an error occurs when parsing, then None is returned
        """
        try:
            inst = cls.from_json(data)
            return cls if cls.is_valid_instance(inst) else None
        except (KeyError):
            return None

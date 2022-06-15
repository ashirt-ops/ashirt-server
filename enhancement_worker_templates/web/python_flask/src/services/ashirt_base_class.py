from abc import ABC, abstractmethod
from base64 import b64decode
import json
from typing import Literal, Optional

from . import (
    make_hmac,
    now_in_rfc1123,
    RequestConfig as RC,
    CreateOperationInput,
)


class AShirtService(ABC):
    """
    AShirtService is an abstract class that holds the necessary details to construct a request with
    the proper headers to contact the AShirt backend. Note that this goes up to modeling the request.
    The actual sending of the request is left to the subclasses.
    """
    def __init__(self, api_url: str, access_key: str, secret_key_b64: str):
        self.api_url = api_url
        self.access_key = access_key
        self.secret_key = b64decode(secret_key_b64)

    @abstractmethod
    def _make_request(cls, cfg: RC, headers: dict[str, str], body: Optional[bytes]):
        """
        _make_request is an abstract method designed to actually make the request. Subclasses will
        need to implement this with the boilerplate code that actually does the request.
        """
        pass

    def get_operations(self):
        return self.build_request(RC('GET', '/api/operations'))

    def create_operation(self, i: CreateOperationInput):
        return self.build_request(RC('POST', '/api/operations', json.dumps(i)))

    def check_connection(self):
        return self.build_request(RC('GET', '/api/checkconnection'))

    def get_evidence_content(self, operation_slug: str, evidence_uuid: str, content_type: Literal['media', 'preview']):
        return self.build_request(RC(
            'GET',
            f'/api/operations/{operation_slug}/evidence/{evidence_uuid}/{content_type}',
            None,
            'raw'
        ))

    def build_request(self, cfg: RC):
        """
        build_request models a request, and the passes the request to the actual executor methods
        (_make_request)
        """
        now = now_in_rfc1123()

        # with_body should now be either bytes or None
        with_body = cfg.body.encode() if type(cfg.body) == str else cfg.body

        auth = make_hmac(cfg.method, cfg.path, now, with_body,
                         self.access_key, self.secret_key)
        headers = {
            "Content-Type": "application/json",
            "Date": now,
            "Authorization": auth,
        }

        return self._make_request(cfg, headers, with_body)

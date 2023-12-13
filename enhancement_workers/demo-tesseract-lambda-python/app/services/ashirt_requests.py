from base64 import b64decode
import json
from typing import Optional, Literal

import requests

from . import (
    make_hmac,
    now_in_rfc1123,
    RequestConfig as RC,
    CreateOperationInput,
)


class AShirtRequestsService:
    def __init__(self, api_url: str, access_key: str, secret_key_b64: str):
        self.api_url = api_url
        self.access_key = access_key
        self.secret_key = b64decode(secret_key_b64)

    ### Request methods to AShirt
    def get_operations(self):
        return self.build_request(RC('GET', '/api/operations'))

    def create_operation(self, i: CreateOperationInput):
        return self.build_request(RC('POST', '/api/operations', json.dumps(i)))

    def check_connection(self):
        return self.build_request(RC('GET', '/api/checkconnection'))

    def get_evidence_content(self, operation_slug: str, evidence_uuid: str, content_type: Literal['media', 'preview']='media'):
        return self.build_request(RC(
            'GET',
            f'/api/operations/{operation_slug}/evidence/{evidence_uuid}/{content_type}',
            None,
            'raw'
        ))

    ### Request helpers

    def build_request(self, cfg: RC):
        """
        build_request models a request, and the passes the request to the actual executor methods
        (_make_request)
        """
        now = now_in_rfc1123()

        # with_body should now be either bytes or None
        with_body = cfg.body.encode() if type(cfg.body) is str else cfg.body

        auth = make_hmac(cfg.method, cfg.path, now, with_body,
                         self.access_key, self.secret_key)
        headers = {
            "Content-Type": "application/json",
            "Date": now,
            "Authorization": auth,
        }

        return self._make_request(cfg, headers, with_body)

    def _make_request(self, cfg: RC, headers: dict[str, str], body: Optional[bytes]):
        resp = requests.request(
            cfg.method, f'{self.api_url}{cfg.path}', headers=headers, data=body)
        if cfg.return_type == 'json':
            return resp.json()
        return resp.content

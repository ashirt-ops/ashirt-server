from typing import Optional
import requests

from .ashirt_base_class import AShirtService
from . import (
    RequestConfig as RC,
)


class AShirtRequestsService(AShirtService):
    """
    AShirtRequestsService is a subclass of AShirtService that makes requests using the Requests
    library. This is a sychronous library, and so care needs to be taken when using this service.
    """
    def __init__(self, api_url: str, access_key: str, secret_key_b64: str):
        super().__init__(api_url, access_key, secret_key_b64)

    def _make_request(self, cfg: RC, headers: dict[str, str], body: Optional[bytes])->bytes:
        resp = requests.request(
            cfg.method, self._route_to(cfg.path), headers=headers, data=body, stream=True)

        if cfg.return_type == 'json':
            return resp.json()
        elif cfg.return_type == 'status':
            return resp.status_code
        elif cfg.return_type == 'text':
            return resp.text

        return resp.content

    def _route_to(self, path: str):
        return f'{self.api_url}{path}'

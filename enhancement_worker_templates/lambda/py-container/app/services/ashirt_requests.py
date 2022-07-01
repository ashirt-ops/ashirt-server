from base64 import b64decode
import json
from typing import Any, Callable, Literal, Optional

import requests

from . import (
    encode_form,
    make_hmac,
    now_in_rfc1123,
    RequestConfig as RC,
    CreateOperationInput,
    CreateEvidenceInput,
    CreateTagInput,
    UpdateEvidenceInput,
    UpsertEvidenceMetadata,
    ListOperationsOutput,
    OperationOutputItem,
    CheckConnectionOutput,
    ReadEvidenceOutput,
    EvidenceOutput,
    ListOperationTagsOutput,
    TagOutputItem,
)


class AShirtRequestsService:
    def __init__(self, api_url: str, access_key: str, secret_key_b64: str):
        self.api_url = api_url
        self.access_key = access_key
        self.secret_key = b64decode(secret_key_b64)

    ### Request methods to AShirt
    def get_operations(self) -> ListOperationsOutput:
        return self.build_request(RC('GET', '/api/operations'))

    def create_operation(self, i: CreateOperationInput) -> OperationOutputItem:
        return self.build_request(RC('POST', '/api/operations', json.dumps(i)))

    def check_connection(self) -> CheckConnectionOutput:
        return self.build_request(RC('GET', '/api/checkconnection'))

    def get_evidence(self, operation_slug: str, evidence_uuid: str) -> ReadEvidenceOutput:
        return self.build_request(RC('GET', f'/api/operations/{operation_slug}/evidence/{evidence_uuid}'))

    def get_evidence_content(self, operation_slug: str, evidence_uuid: str, content_type: Literal['media', 'preview']='media'):
        return self.build_request(RC(
            'GET',
            f'/api/operations/{operation_slug}/evidence/{evidence_uuid}/{content_type}',
            None,
            'raw'
        ))

    def create_evidence(self, operation_slug: str, i: CreateEvidenceInput) -> EvidenceOutput:
        body = {
            'notes': i['notes'],
        }
        add_if_not_none(body, 'contentType', i.get('content_type'))
        add_if_not_none(body, 'tagIds', i.get('tag_ids'), json.dumps)

        data = encode_form(body, {"file": i.get('file')})

        return self.build_request(RC('POST',
            f'/api/operations/{operation_slug}/evidence',
            body=data['data'],
            multipart_boundary=data['boundary'])
            )

    def update_evidence(self, operation_slug: str, evidence_uuid: str, i: UpdateEvidenceInput) -> int:
        """update_evidence revises evidence per the given input. Returns a 200 response if successful."""
        body = {}

        add_if_not_none(body, 'notes', i.get('notes'))
        add_if_not_none(body, 'contentType', i.get('content_type'))
        add_if_not_none(body, 'tagsToAdd', i.get('add_tag_ids'), json.dumps)
        add_if_not_none(body, 'tagsToRemove', i.get('remove_tag_ids'), json.dumps)

        data = encode_form(body, {"file": i.get('file')})

        return self.build_request(RC('PUT',
            f'/api/operations/{operation_slug}/evidence/{evidence_uuid}',
            body=data['data'],
            multipart_boundary=data['boundary'],
            return_type='status'
        ))

    def upsert_evidence_metadata(self, operation_slug: str, evidence_uuid: str, i: UpsertEvidenceMetadata):
        """upsert_evidence_metadata sets or updates evidence metadata per the given input.
        Returns a 200 response if successful.
        """
        return self.build_request(RC(
            'PUT',
            f'/api/operations/{operation_slug}/evidence/{evidence_uuid}/metadata',
            body=json.dumps(i),
            return_type='status'
        ))

    def get_operation_tags(self, operation_slug: str) -> ListOperationTagsOutput:
        return self.build_request(RC('GET', f'/api/operations/{operation_slug}/tags'))

    def create_operation_tag(self, operation_slug: str, i: CreateTagInput) -> TagOutputItem:
        return self.build_request(RC('POST', f'/api/operations/{operation_slug}/tags', json.dumps(i)))

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

        if cfg.multipart_boundary is None:
            content_type = "application/json"
        else:
            content_type = f'multipart/form-data; boundary={cfg.multipart_boundary}'

        headers = {
            "Content-Type": content_type,
            "Date": now,
            "Authorization": auth,
        }

        return self._make_request(cfg, headers, with_body)

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


def add_if_not_none(body: dict[str, Any], key: str, value: Any, tf: Callable[[Any], Any]=None):
    if value is not None:
        body.update({key: value if tf is None else tf(value)})

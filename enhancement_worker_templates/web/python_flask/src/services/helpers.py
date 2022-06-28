from base64 import b64encode, urlsafe_b64encode
from datetime import datetime
import hashlib
import hmac
import os
from typing import Optional
from wsgiref.handlers import format_date_time

from .types import HTTP_METHOD, FileData, MultipartData


def make_hmac(
    method: HTTP_METHOD,
    path: str,
    date: str,
    body: Optional[bytes],
    access_key: str,
    secret_key: bytes
):
    """
    make_hamc builds the authentication string needed to contact ashirt.
    """
    body_digest_method = hashlib.sha256()
    if body is not None:
        body_digest_method.update(body)
    body_digest = body_digest_method.digest()

    to_be_hashed = f'{method}\n{path}\n{date}\n'
    full_message = to_be_hashed.encode() + body_digest

    hmacMessage = b64encode(
        hmac.new(secret_key, full_message, hashlib.sha256).digest())

    return f'{access_key}:{hmacMessage.decode("ascii")}'


def now_in_rfc1123():
    """now_in_rfc1123 constructs a date like: Wed, May 11 2022 09:29:02 GMT"""
    return format_date_time(datetime.now().timestamp())


def _random_char(length: int):
    return urlsafe_b64encode(os.urandom(length))


def encode_form(fields: dict[str, str], files: dict[str, FileData]) -> MultipartData:
    boundary = "----AShirtFormData-".encode() + _random_char(30)
    newline = "\r\n".encode()
    part = "--".encode()
    boundary_start = part + boundary + newline
    last_boundary = part + boundary + part + newline
    content_dispo = "Content-Disposition: form-data".encode()

    field_buff = bytes()
    for key, value in fields.items():
        entry = (
            boundary_start +
            content_dispo + f'; name="{key}"'.encode() +
            newline + newline +
            value.encode() +
            newline
        )
        field_buff += entry

    file_buff = bytes()
    for key, value in files.items():
        if value is None:
            continue
        entry = (
            boundary_start +
            content_dispo + f'; name="{key}"; filename="{value["filename"]}"'.encode() +
            newline + f'Content-Type: {value["mimetype"]}'.encode() +
            newline + newline +
            value['content'] +
            newline
        )
        file_buff += entry

    return {
        "boundary": boundary.decode(),
        "data": field_buff + file_buff + last_boundary
    }

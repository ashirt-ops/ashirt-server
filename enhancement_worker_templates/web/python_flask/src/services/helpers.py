from base64 import b64encode
import binascii
from datetime import datetime
import hashlib
import hmac
import os
from typing import Optional
from wsgiref.handlers import format_date_time

from .types import HTTP_METHOD, FileData


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
    return binascii.hexlify(os.urandom(length))


def encode_form(fields: dict[str, str], files:dict[str, FileData]):
    boundry = "----AShirtFormData-".encode() + _random_char(30)
    newline = "\r\n".encode()
    boundry_start = boundry + newline
    last_boundry = boundry + "--".encode() + newline

    body = bytes()

    def field_part(name: str, value: str | bytes, extra=""):
        header = f'Content-Disposition: form-data; name="{name}"{extra}\r\n'.encode(
        )
        return (
            boundry_start +
            header +
            newline +
            (value.encode() if type(value) is str else value)
        )

    def file_part(name, value, filename, content_type):
        extra = f'; filename="{filename}"\r\nContent-Type: {content_type}'
        return field_part(name, value, extra)

    for k, v in fields.items():
        body += field_part(k, v)
        body += newline

    for k, v in files.items():
        body += file_part(k, v['content'], v['filename'], v['mimetype'])
        body += newline

    body += last_boundry

    return body

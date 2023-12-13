from base64 import b64encode
from datetime import datetime
import hashlib
import hmac
from typing import Optional
from wsgiref.handlers import format_date_time

from .types import HTTP_METHOD


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

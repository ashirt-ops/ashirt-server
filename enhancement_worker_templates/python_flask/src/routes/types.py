from enum import Enum


class StatusCode(Enum):
    """StatusCode is a set of status codes supported by AShirt."""
    OK = 200
    ACCEPTED = 202
    NO_CONTENT = 204
    BAD_REQUEST = 400
    NOT_ACCEPTABLE = 406
    INTERNAL_SERVICE_ERROR = 500
    NOT_IMPLEMENTED = 501

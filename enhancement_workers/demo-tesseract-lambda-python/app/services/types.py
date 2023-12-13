from dataclasses import dataclass
from typing import Literal, Optional, TypedDict, Union

HTTP_METHOD = Literal['GET', 'POST', 'PUT', 'DELETE']


@dataclass(frozen=True)
class RequestConfig:
    """
    RequestConfig abstracts a request so that it can be sent via different libraries,
    in case you don't like requests
    """
    method: HTTP_METHOD
    path: str
    body: Optional[Union[bytes, str]] = None
    return_type: Literal["json", "raw"] = "json"


# The below are all inputs for various API calls

class CreateOperationInput(TypedDict):
    slug: str
    name: str

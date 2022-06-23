from typing import Literal, Optional, TypedDict, Union


class ProcessResultNormal(TypedDict):
    action: Literal['rejected', 'error']
    content: Optional[str]


class ProcessResultComplete(TypedDict):
    action: Literal['processed']
    content: str


class ProcessResultDeferred(TypedDict):
    action: Literal['deferred']


ProcessResultDTO = Union[ProcessResultNormal,
                         ProcessResultComplete,
                         ProcessResultDeferred]

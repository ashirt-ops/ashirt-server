from .helpers import *
from .types import *
from .ashirt_base_class import *
from .ashirt_sync import *


_ashirt_service: AShirtService


def set_service(svc: AShirtService):
    """
    set_service stores an instance of a concrete AShirtService class. This is paired with svc to
    allow making requests anywhere in the appliction.
    """
    global _ashirt_service
    _ashirt_service = svc


def svc() -> AShirtService:
    """
    svc provides an established, concrete AShirtService class that can make requests to AShirt.
    """
    return _ashirt_service

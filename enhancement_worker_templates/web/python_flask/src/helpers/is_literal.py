from typing import Any


def is_literal(v: Any, expectedType: type, expectedValue: Any) -> bool:
    """
    is_literal is a small helper that verifies that the value passed has the expected type
    and the expected value. This is useful to validate literal values provided by an external
    service
    """
    return (
        type(v) == expectedType
        and v == expectedValue
    )

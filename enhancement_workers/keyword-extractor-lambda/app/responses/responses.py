import json
from typing import Optional

def test_passed():
    return {
        'statusCode': 200,
        'body': json.dumps({
            'status': 'ok'
        })
    }

def process_success(content):
    return {
        'statusCode': 200,
        'body': json.dumps({
            'action': 'processed',
            'content': content
        })
    }

def reject_evidence():
    return {
        'statusCode': 406,
    }


def not_implemented():
    return {
        'statusCode': 501,
    }


def error_processing(message: Optional[str] = None):
    return _error_resp(500, message)


def bad_request(message: Optional[str] = None):
    return _error_resp(400, message)


def _error_resp(status_code: int, message: Optional[str] = None):
    rtn = {
        'statusCode': status_code,
    }

    if message is not None:
        rtn['body'] = json.dumps({
            'action': 'error',
            'content': message
        })

    return rtn

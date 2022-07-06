from flask import (
    Blueprint, request, current_app, jsonify, Response, g
)
import json
from uuid import uuid4

from constants import APP_LOGGER
from request_types import (EvidenceCreatedBody, TestBody)
from state import RequestState
import actions

from .types import StatusCode


bp = Blueprint('ashirt', __name__, url_prefix='/ashirt')


@bp.route("/process", methods=['POST'])
def process_request() -> Response:
    """
    process_request handles requests received from AShirt
    """
    data = request.json
    if TestBody.parse_if_valid(data) is not None:
        return jsonify({"status": "ok"})

    if (body := EvidenceCreatedBody.parse_if_valid(data)) is not None:
        action_result = actions.handle_evidence_created(body)
        # Construct a response that provides a body when a body is meaningful
        rtn = (
            Response()
            if action_result.get('content') is None
            else Response(json.dumps(action_result))
        )

        rtn.status_code = {
            'processed': StatusCode.OK.value,
            'deferred': StatusCode.ACCEPTED.value,
            'error': StatusCode.INTERNAL_SERVICE_ERROR.value,
            'rejected': StatusCode.NOT_ACCEPTABLE.value,
        }[action_result['action']]
        return rtn

    return Response('Unsupported Body Type', status=501)


########## Blueprint Stuff ############

@bp.before_request
def on_request_received():
    """
    on_request_received established a state for the request, complete with a logger. Also logs the
    start of the request for anaylitics purposes
    """
    ctx = str(uuid4())
    app_logger = current_app.config[APP_LOGGER]
    req_log = app_logger.bind(context=ctx)

    req_log.msg("Received Request",
                method=request.method,
                endpoint=request.full_path,
                query=request.query_string)
    g._request_state = RequestState(req_log)


@bp.after_request
def on_request_complete(resp):
    """
    on_request_complete logs when the request has been completed
    """
    req_state = g._request_state
    if type(req_state) == RequestState:
        g._request_state.req_log.msg(
            "Request Complete", response_code=resp.status_code, body=resp.data)

    return resp

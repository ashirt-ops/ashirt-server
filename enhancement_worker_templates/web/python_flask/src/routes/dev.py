from flask import (
    Blueprint, jsonify, Response
)

# from services import svc


bp = Blueprint('', __name__, url_prefix='/')


@bp.route("/")
def index() -> Response:
    """index provides a method to verify that the service is live"""
    return jsonify({
        "msg": "GET /"
    })

@bp.route("/test", methods=['POST'])
def test() -> Response:
    """test provides a place to verify that individual steps work as expected"""

    return jsonify({"Done": "you bet!"})

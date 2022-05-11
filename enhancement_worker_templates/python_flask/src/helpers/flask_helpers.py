from flask import (
    current_app, Response, make_response, Flask,
)


def jsonify_no_content() -> Response:
    """
    jsonify_no_content produces a 204 (no content) response
    """
    # from https://www.erol.si/2018/03/flask-return-204-no-content-response/
    response = make_response('', 204)
    response.mimetype = current_app.config['JSONIFY_MIMETYPE']
 
    return response


def remove_flask_logging(app: Flask) -> None:
    # See: https://gist.github.com/daryltucker/e40c59a267ea75db12b1
    import logging
    app.logger.disabled = True
    log = logging.getLogger('werkzeug')
    log.disabled = True

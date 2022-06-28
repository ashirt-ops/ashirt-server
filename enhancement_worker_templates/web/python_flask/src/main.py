import os

from dotenv import dotenv_values
from flask import Flask
import structlog

from constants import APP_LOGGER, STATE_NAME
from helpers import remove_flask_logging
from project_config import ProjectConfig
from routes import (ashirt, dev)
from services import AShirtRequestsService
from services import set_service


def create_app() -> Flask:
    app = Flask(__name__)

    full_env = {
        **dotenv_values(".env"),
        **os.environ
    }
    cfg = ProjectConfig.from_dict(full_env)
    app.config[STATE_NAME] = cfg
    app.config[APP_LOGGER] = structlog.get_logger()

    set_service(
        AShirtRequestsService(cfg.backend_url, cfg.access_key, cfg.secret_key_b64)
    )
    app.register_blueprint(ashirt.bp) # Add normal routes
    if cfg.dev_mode:
        app.config[APP_LOGGER].msg("Adding dev routes")
        app.register_blueprint(dev.bp) # Add dev routes

    # tweak logging settings
    remove_flask_logging(app)
    return app


if __name__ == "__main__":
    app = create_app()
    try:
        app.config[APP_LOGGER].msg("App Starting")
        app.run(host="0.0.0.0", port=app.config[STATE_NAME].port)
    finally:
        app.config[APP_LOGGER].msg("App Exiting")

# Python/Flask Enhancement Worker

This worker is a microservice-type worker built on Python/Flask.
It uses [pipenv](https://pipenv.pypa.io/en/latest/) to manage dependencies.
This project strives to be as minimalistic as possible, but does include some helpful libraries. This
includes:

* [Flask](https://flask.palletsprojects.com/en/2.1.x/), to manage the network connection
* [gunicorn](https://gunicorn.org/), for production deployment
* [requests](https://docs.python-requests.org/en/latest/), to handle contacting the ashirt instance
* [structlog](https://www.structlog.org/en/stable/), for structured logging
* [python-dotenv](https://pypi.org/project/python-dotenv/), for environment loading (this is primarily aimed at development)

In addition, this service tries to be as type-safe as possible, so extra effort has been provided to ensure that the typing is specified as much as possible.

To get up and running, open the project root in a terminal, install pipenv, and run `pipenv shell`, then `pipenv install`

## Deploying to AShirt

The typical configuration for deploying this worker archetype is going to look roughly like this:

```json
{
    "type": "web", 
    "version": 1,
    "url": "http://testapp/ashirt/process"
}
```

Note the url: this is likely what will change for your version.

## Adding custom logic

Most programs should be able to largely ignore most of the code, and instead focus on `actions` directory, and specifically the events you want to target.

## How it works

This section is mostly for those that need to do more than implement the core functionality. This application, like many other webservices, can be divided up into two phases: the startup phase, and the serving phase.

### Startup Phase

The Startup Phase is as you might expect: this state is entered once the application starts, and it is responsible for configuring the application for long-term running. The most important bit here is likely the configuration and route management. `create_app` within `main.py` will load configuration details from the environment (locally: `.env` file), create a class for handling requests to an AShirt backend, and register standard routes. Then, either the main line in `main.py` or `wsgi.py` will start the server. This phase ends once the server starts servicing requests, and the application then enters the serving phase.

### Serving Phase

The serving phase is largely controlled by what particular route is entered when a user contacts the server. The `routes` directory provides two set of routes: the `ashirt` routes, which are the routes that the AShirt backend will call, and the `dev` routes, which are designed to be only created in a development environment. These serve as helpers and sanity checks.

When a route is reached -- in particular, when the process route is reached (see: `process_request`), then the service will kick off processing of that data. Some initial boilerplate style code manages the request, and directs all of the actual work to some function in the `actions` directory. These functions will return one of a handful of responses, which will be used to generate the true response to the AShirt backend.

Once the request is complete, the application waits for another request.

### Contacting AShirt

The `services` folder contains a class that is used to contact ashirt. This is the typical way of getting the actual content / interacting with ashirt. This is treated mostly as a singleton by preparing the service in `main.py`, and recording an instance in `services/__init__.py`. Any other module can then use the `svc` function to get the loaded instance.

## Integrating into AShirt testing environment

Notably, the dev port exposed is port 8080, so all port mapping has to be done with that in mind. When running locally (not via docker), the exposed port is configurable.

This configuration should work for your scenario, though the volumes mapped might need to be different.

```yaml
  py_flask:
    build:
      context: enhancement_worker_templates/web/python_flask
      dockerfile: Dockerfile.dev
    ports:
      - 3004:8080
    restart: on-failure
    volumes:
      - ./enhancement_worker_templates/web/python_flask/src:/app/src
    environment:
      ENABLE_DEV: true
      ASHIRT_BACKEND_URL: http://backend:3000
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```


Note that the mapped volume overwrites the source files placed in the image. This allows for hot-reloading of the worker when deployed to docker-compose. If you don't want or need hot reloading, then you can simply omit this declaration.

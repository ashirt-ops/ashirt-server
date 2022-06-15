# Python/Flask Enhancement Worker

This worker is a microservice-type worker built on Python/Flask.
It uses [pipenv](https://pipenv.pypa.io/en/latest/) to manage dependencies.
This work strives to be as minimalistic as possible, but does include some creature comforts. This
includes:

* [Flask](https://flask.palletsprojects.com/en/2.1.x/), to manage the network connection
* [gunicorn](https://gunicorn.org/), for production deployment
* [requests](https://docs.python-requests.org/en/latest/), to handle contacting the ashirt instance
* [structlog](https://www.structlog.org/en/stable/), for structured logging
* [python-dotenv](https://pypi.org/project/python-dotenv/), for environment loading (this is primarily aimed at development)

In addition, this service tries to be as type-safe as possible, so extra effort has been provided to ensure that the typing is specified as much as possible.

To get up and running, open the project root in a terminal, install pipenv, and run `pipenv shell`, then `pipenv install`

## Configuration

The typical configuration for deploying this worker archetype is going to look roughly like this:

```json
{
    "type": "web", 
    "version": 1,
    "url": "http://testapp/ashirt/process"
}
```

Note the url: this is likely what will change for your version

## Adding custom logic

Most programs should be able to largely ignore most of the code, and instead focus on `actions/process_handler.py`. The `handle_process` function is ultimately called when new evidence is added to AShirt. Simply add in your logic here to process new pieces of evidence, and you should be good to go.

## How it works

You only need this if you need more than simply adding your core functionality. This application, like many other webservices, can be divided up into two phases: the startup phase, and the serving phase.

### Startup Phase

The Startup Phase is as you might expect: this state is entered once the application starts, and is responsible for configuring the application for long-term running. The most important bit here is likely the configuration and route management. `create_app` within `main.py` will load configuration details from the environment (locally: `.env` file), create a class for handling requests to an AShirt backend, and register standard routes. Then, either the main line in `main.py` or `wsgi.py` will start the server. This phase ends once the server starts servicing requests, and the application then enters the serving phase.

### Serving Phase

The serving phase is largely controlled by what particular route is entered when a user contacts the server. The `routes` directory provides two set of routes: the `ashirt` routes, which are the routes that the AShirt backend will call, and the `dev` routes, which are designed to be only created in a development environment. These serve as helpers and sanity checks.

When a route is reached -- in particular, when the process route is reached (see: `process_request`), then the service will kick off processing of that data. Some initial boilerplate style code manages the request, and directs all of the actual work to the `process_handler.py` file. This function will return one of a handful of responses, which will be used to generate the true response to the AShirt backend.

Once the request is complete, the application waits for another request.

### Contacting AShirt

The `services` folder contains a class that is used to contact ashirt. This is the typical way of getting the actual content / interacting with ashirt. This is treated mostly as a singleton by preparing the service in `main.py`, and recording an instance in `services/__init__.py`. Any other module can then use the `scv` function to get the loaded instance.

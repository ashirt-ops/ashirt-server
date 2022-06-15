# Typescript/Express Enhancement Worker

This worker is a microservice-type worker built on Typescriipt/Express.
It uses [yarn](https://yarnpkg.com/) to manage dependencies, though [npm](https://www.npmjs.com/) works just as well. This project strives to be as minimalistic as possible, but does include some creature comfots. This includes:

* [Express](https://expressjs.com/), to handle receiving network requests
* [Typescript](https://www.typescriptlang.org/) As the language of choice
* [axios](https://axios-http.com/), to make network connections
* [dotenv](https://github.com/motdotla/dotenv#readme), to load environment variables

## Configuration

The typical configuration for deploying this worker archetype is going to look roughly like this:

```json
{
    "type": "web", 
    "version": 1,
    "url": "http://testapp/process"
}
```

Note the url: this is likely what will change for your version

## Adding custom logic

Most applications can largely ignore the majority of the application and instead focus on just the file at `src/actions/processAction.ts`, and the `handleActionProcess` function in particular. This function is ultimately called when a process request arrives and passes basic validation. From this function, you can focus completely on how a recevied request should be handled.

## How it works

You only need to understand this if you need more than simply adding your core functionality. This application, like many other webservices, can be divided up into two phases: the startup phase, and the serving phase.

### Startup Phase

The Startup Phase is as you expect: entered once the application starts, and is responsible for configuring the application for long-term running. The most important bit here is likely the configuration and route management. `src/main.ts` bootstaps the application, loading the configuration (defined in `src/config.ts`), then defining the standard routes (defined in `src/router.ts`) before ultimately starting the service. This phase ends once the server starts servicing requests, and the application then enters the serving phase.

### Serving Phase

The serving phase is largely controlled by what particular route is entered when a user contacts the server. The `src/router.ts` file provides two types of endpoints. The standard `process` endpoint to handle AShirt requests, and a set of dev routes, which are only enabled when in dev mode. The dev routes act as sanity checks and direct access when testing new features.

When a route is reached -- in particular, when the process route is reached (see: `app.post('/process'`), then the service will kick off processing of that data. Some initial boilerplate style code manages the request, deferring all of the actual work to the `src/actions.processAction.ts` file. This function will return one of a handful of responses, which will be used to generate the true response to the AShirt backend.

Once the request is complete, the application waits for another request.

### Contacting AShirt

The `services` folder contains a class that is used to contact ashirt. This is the typical way of getting the actual content / interacting with ashirt. The instance of this class is available to the routes defined in `addRoutes`, which can then be passed into the handlers as needed.

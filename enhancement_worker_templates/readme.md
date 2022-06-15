# Enhancement Worker Templates

This directory acts as a repository for service-based workers. A demo worker for testing can be found in enhancement_workers. You can use this as a jumping off place for your own worker.

Each of the template workers contained in here have their own readmes describing how the template works. The goal for these templates is to provide a reasonable starting point with much of the boiler plate code written and tested. Currently, the workers fall into two big categories:

* Web-based workers
* AWS Lambda-based workers

Templates are thus broken up as such:

* `web/` All of the web-based workers
* `lambda/` All of the AWS Lambda based workers

The projects themselves are then broken up by what language and/or framework they use. Feel free to look among them, or build your own.

While these templates have their own concerns and patterns, they share a lot. The rest of this document attempts to highlight that commonality.

## Adding a service to AShirt

All AShirt enhancement workers have a configuration that is specified to AShirt, which provides the details on how to contact the service. The workers should all have the configuration documented, but each looks like the following:

```json
{
    "type": "aws", // what broad category of service this falls into. Currently we support "web" and "aws" (lambda)
    "version": 1, // A version number for the config, in case additional configurations are required.
    // and additional, type-sepcific configurations are also contained here
}
```

Get this configuration, tweak it if necessary, and as an admin, create a new worker by specifying this config to AShirt. Go to `/admin/services` and click the "Create  New Service Worker" button to create the service.

## Using a service for development purposes

At certain points, you may want to incorporate a worker into your local AShirt development environment. This is pretty straight forward for web and container-based lambda workers. Layer-based lambda workers may need to be converted into their container bretheren before they are testable. In both of these cases do the following:

1. Open the `docker-compose.yml` file.
2. Somewhere under `services`, add the following config:

```yml
  myWorker:
    build:
      context: path/to/worker/myWorker
      dockerfile: Dockerfile
    ports:
      - 3001:3001
    restart: on-failure
    environment:
      PORT: 3001
      ENABLE_DEV: true
      ASHIRT_BACKEND_URL: http://backend:3000
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```

A lot of the values may change, depending on your needs. Here is what you need to know:

* `[worker-name]` (the top most yaml declaration): This name can be anything, but is important. Your configuration will need to be tweaked based on this value so you can properly contact it. For example, if you have a webservice, the url config should look like: `"url": "http://[worker-name]:3001/process",`. If using a lambda worker, then the `[worker-name]` will need to match the `lambdaName` configuration value.
* `Ports`: the particular port mapping doesn't matter, but will be specified in the config. Note that for container-based lambda workers, the exposed port is 8080, so your mapping should look like `3001:8080`. Here, `3001`  was chosen as the first port close to the backend port that was not already occupied.
* `ASHIRT_ACCESS_KEY` and `ASHIRT_SECRET_KEY` need not be changed if you are using the normal ashirt seeding. A headless user with these keys will be added to the database upon creation.

All other values are free to edit, but assume you understand how to configure your application, and how to specify the configuration to docker-compose.

### Permanently loading into AShirt seeding

If you want to keep this service around for awhile, you can define the docker-compose value above, and in addition, you can add the configuration and service name to the `hp_seed_data.go` (search for `newHPServiceWorker` to find examples on how to create your own)

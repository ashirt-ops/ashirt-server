# Lambda Template - Container/Javascript

This template creates a lambda worker via a [docker image base](https://docs.aws.amazon.com/lambda/latest/dg/images-create.html). This is generally the preferred method when creating a lambda for AShirt, as it makes testing offline a bit easier.

## Editing

As with the standard layered javascript template, the boiler plate is largely complete. The majority of the work can be focused on filling out `handleEvidenceCreated` and `doProcessing` in app.js. The rest of the code attempts to be as dependency free as possible.

If you need to add dependencies, you can create a `package.json` file to define these, as you would with any nodejs project. The docker base image can likewise be expanded to both install the dependencies as well as add any extra dependencies/software that is needed.

## Deploying to AWS

AWS provides documentation on how to deploy these functions [here](https://docs.aws.amazon.com/lambda/latest/dg/images-create.html#images-upload)

## Deploying to AShirt

The standard AShirt configuration can be used with these types of workers:

```ts
{
    type: 'aws',
    version: 1,
    lambdaName: string
    asyncFunction: bool
}
```

### Adding a test deployment to AShirt's docker compose file

The following docker-compose configuration section should work:

```yml
  myWorker:
    build:
      context: path/to/project
      dockerfile: Dockerfile
    ports:
      - 3001:8080
    restart: on-failure
    environment:
      ASHIRT_BACKEND_URL: backend
      ASHIRT_BACKEND_PORT: 3000
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```

Here, update the worker name (must match the lambdaName in the configuration), and you can optionally update the `3001` portion of the `ports` configuration. Likewise, you can add more values to environment as needed. Note that `ASHIRT_ACCESS_KEY` and `ASHIRT_SECRET_KEY` can be kept as is for test deployments using the standard seed, as these values are already baked into the seeding.

## Testing

A small makefile has been created to provide baseline building and testing. This is focused purely on if this image can be successfully built and executed without encountering errors. Feel free to use this to do some initial verification, but you may also want to look into unit and integration tests as your scenarios get more complex.

To build, run `make build`, to run, execute `make run`, or to do both, use `make start`. To test, in a separate terminal, run `make test-test` to send a standard test command, or `make test-process`, which will send a process-type request. You can also test bad scenarios, if needed, via `make test-unsupported`

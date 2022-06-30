# Lambda Template - Layer/Javascript

This template creates a lambda worker via a [layer](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-concepts.html#gettingstarted-concepts-layer). While this works, it is recommended that lambdas be based containers for better maintainability.

## Editing

As with the standard layered javascript template, the boiler plate is largely complete. The majority of the work can be focused on filling out `handleEvidenceCreated` and `doProcessing` in app.js. The rest of the code attempts to be as dependency free as possible.

If you need to add dependencies, you can create a `package.json` file to define these, as you would with any nodejs project. These dependencies cannot be carried over in a traditional way, and will instead need to be transformed into layers and uploaded to aws lambda before you can use them.

## Deploying to AWS

AWS contains documentation on how to deploy nodejs layered lambda functions [here](https://docs.aws.amazon.com/lambda/latest/dg/nodejs-package.html)

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

Unfortunately, the layer-based lambdas are not easily integration-testable. It is recommended you add a docker image (see [here](/enhancement_worker_templates/lambda/js-container/readme.md) for container-based lambdas) and test that way. However, if you want to test with your deployment, you'll need to actually deploy the function, as well as your backend, so that they may speak with each other. That's beyond the scope of this documentation however.

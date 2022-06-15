# Lambda Worker Templates

AWS Lambda workers are enhacement workers that execute on aws lambda. This provides for a more streamlined deployment experience, but may not be right for every worker.

Lambda provides two main types of lambda deployments: layer based (the original way) and container based. In general, these are simply packaging concerns, and the code written for one way is largely applicable to the other. That said, it is generally easier to create a container for local testing, so that is the recommended solution.

Workers in this directory are named by the language they use, and the deployment mechanism they use. For example `js-container` refers to the JavaScript programming language, for a container-based deployment. js-layer likewise refers to a layer-based deployment.

## AShirt Config / deployment

All Lambda configs will have the following configuration syntax:

```ts
{
    type: 'aws',
    version: 1,
    lambdaName: string
    asyncFunction: bool
}
```

When deploying for testing, keep in mind that the lambda image expects network calls on port 8080, and so you must map to that port.

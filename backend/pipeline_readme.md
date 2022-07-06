# AShirt API Event System

The AShirt API Event system enables enrichment and automation opportunties for your ashirt deployment. While this is primarily aimed at providing metadata for evidence, each worker in this sytem communicates via the API, and so has access to all API methods.

Additionally, while AShirt provides some enrichment services, the definition is kept open so that you may create your own services, or use community-created services.

## A Word on Security

Please review and consider each new service you add to your AShirt instance. These services get direct access to your evidence, and via API, have access to a wide slice of your AShirt data. As such, it is important that any service you add here has been vetted by you and your team.

## Installing/Adding a Service

To add a service worker, as an admin, navigate to admin/service workers (url: `/admin/services`). Click the "Create New Service Worker" button and specify a name for the worker, as well as the configuration. The configuration will be a JSON body, typically provided by the service itself. Generally this will contain details on how to contact the service. See [Below](#web-version-1) for the web configuration schema as an example of what to place here. Once the name and configuration have been specified, click "Create", and processing of _new evidence_ should begin. Old evidence can be processed on an item-by-item basis, or in a batched manner on an operation-by-operation basis.

## Building a Compliant Service

There are various options in how you can construct a valid service. The primary concern for the AShirt backend is how to contact this service. Once contacted, the service is then responsible for the following:

* Determining if the event is appropriate to process

Additionally, if this is for evidence metadata population:

* Determining if the evidence is appropriate to analyze
* Responding with the result

The exact manner in which the above is accomplished is largely up to the service, so long as it can be configured using a standard configuration template. The below sections specify how to create a configuration hosted by various means.

### Before You Start / Definitions

The details below try to provide all of the necessary details to get you started in building a new pipeline service. To that end, the documentation here assumes you know a little bit about the following:

* What AShirt is, and broadly what it can do
* JSON format, and building/parsing JSON in the language of your choice
* Basic Typescript type definitions. Specifically, you should have an understanding of:
  * Basic javascript types (`string`, `number`, `undefined`, etc)
  * Literal values (`'a'` means literally `'a'` in the message)
  * Optional types (a key with a `?` suffix means this key/value can be omitted)
  * Union types (`a | b` means either `a` or `b`)
  * Record types (objects, or key/value maps. Key and values are restricted to the indicated types)
  * A cheatsheet exists [here](https://www.typescriptlang.org/static/TypeScript%20Types-4cbf7b9d45dc0ec8d18c6c7a0c516114.png)

Additional details may be necessary for the particular type of service you implment.

### Web Service

Web services are any service that can be contacted via an HTTP request, and respond in kind. The details on exactly how the request is sent and what the body is are below. Note that this can include AWS Lambda services, though there is a more specific tool available if you want to use AWS Lambda.

#### Web Configurations

Web services, like all pipeline services, must define their configuration to AShirt when adding the service. The configuration provided below is mostly concerned with how to contact your service, and provides some ways for ashirt to minimally customize the message your service receives.

##### Web, Version 1

```ts
{
  "type": "web",
  "version": 1,
  "url": string,
  "headers": Record<string, string> | undefined // Optional
}
```

#### Handling a Request

AShirt contacts web services via an HTTP POST request. Your service must be able to receive this request, parse the message, and respond back with an appropriate JSON response.

##### Test Connection

AShirt, and the humans using AShirt, would occasionally like to know that the service is operational. To that end, a test message can be sent. The message takes this form:

```ts
{
    "type": "test"
}
```

The service should respond with a 200/OK message, or a 204/No Content response. If everything is okay, the message body may be empty. If you want to communicate additional details, you can provide a response in the following format:

```ts
{
    "status": "ok" | "error"
    "message": string // Optional; A customized message to show to the user
}
```

#### Process Evidence Created Events / Metadata Enhancement

This is called whenever new evidence is added, or on demand for existing pieces of evidence. Either way, the message will have the following format:

```ts
{
  "type": "evidence_created",
  "evidenceUuid": string,
  "operationSlug": string,
  // the below indciates the content type of the evidence. This can help your tool immediately know if processing is worthwhile
  "contentType": 
    | "http-request-cycle" // These are HAR files detailing a request/response session
    | "terminal-recording" // Terminal Recordings are in asciinema file format. 
    | "codeblock"          // Codeblocks are json files. Their format is detailed below.
    | "event"              // Events are no-file pieces of evidence. They act as a time marker
    | "image"              // Images are typically screenshots. Typically these are in PNG format (though this is not a guarantee)
    | "none"               // These are no-file pieces of evidence, containing only a description
}
```

A worker receiving this message has three options:

* The work can be rejected for some reason. For example, the evidence type can't be processed by this worker
* The worker can complete the desired work, and reply to the result in the same request. This is useful if the processing is quick.
* The worker can "defer" the work, and provide a result once processing has completed (via the AShirt API). This is useful for processing that takes some time to analyze. Note that is the responsibility of the worker to manage any queuing of work.

Each of these responses follows the same structure, but is processed according to the context. See below for how to construct these messages.

##### Responding with a status code only

When responding with a status code, it is important that you leave the body empty. Otherwise, json processing of the body will occur, and may trigger an error.

##### Responding by Rejecting

To reject a message, respond with:

* A 406 status code; or
* A 200 status code with a body in the following format:

  ```ts
  {
      "action": "rejected",
      "content": string | undefined, // Optional. Provides an area to explain why the evidence was rejected.
  }
  ```

##### Responding by immediate action

If successful, respond with:

* A 200 response, with the body:

  ```ts
  {
      "action": "processed",
      "content": string, // The result of the processing
  }
  ```

If processing fails, then respond with:

* a 500 status code; or
* A 200 status code with a body in the following format

  ```ts
  {
      "action": "error",
      "content": string | undefined, // Optional. If specified, recorded as the error encountered
  }
  ```

##### Responding by Deferring

To defer work on a message, respond with:

* A 202 status code; or
* A 200 status code with the following body:

  ```ts
  {
      "action": "deferred",
      "content": undefined, // This can be a json string, but it will be ignored
  }
  ```

##### Other Responses

If a response is deliverred that does not match the expected format, the run will be regarded as a failure with a parsing failed error.

See [the API section](#using-the-ashirt-api) on how to contact AShirt once work completes.

### Using the AShirt API

The AShirt API is the medium in which AShirt services and tools can communicate with AShirt and the AShirt database. To communicate, the services need to be attached to a user via an API key and secret. For services, it is recommended a that a headless user is created (this will provide the widest access without having to add a standard user to every operation), and that an API key is generated for that user/service. Once generated, those keys can then be given to the service as a means to construct secure messages to AShirt.

#### Constructing a Message

The vast majority of API calls are JSON requests.

Headers:

* `Content-Type`: `application/json`
* `Date`: now, in `GMT`, in RFC1123 format (note: must be `GMT`, not `UTC`). e.g: `Sun, 21 Oct 2018 12:16:24 GMT` (Also known as RFC7231)
* `Authorization`: HMAC. See below

Authorization is accomplished by constructing an HMAC message. You can find a Golang version in `signer/hmac.go`, in the BuildRequestHMAC function. Likewise, there is a C++ version in AShirt [here](https://github.com/theparanoids/ashirt/blob/main/src/helpers/netman.h#L105-L123) (See the `generateHash` method if the link rusts). However, the process is fairly straight forward, and detailed here:

1. Create the body content. Then, hash this content using the `sha-256` algorithm. Note that in situations where there is no body, you would instead hash an empty string. The output from this should be a series of bytes in no special encoding ("raw" format).
2. Create a string with the following information:
   1. Method (e.g. `GET` or `POST`)
   2. New line
   3. URI/Path of request (e.g. if contacting `http://www.ashirt.com/api/operations`, the path would be `api/operations`)
   4. New line
   5. The `Date` from the headers, defined above
   6. New line
   7. The hashed body, from step 1
3. Take this message, and create an hmac using the `sha-256` algorithm, using your secret key. Convert the result into a base64 string
4. From this, create the Authorization header value in the format: `${API_KEY}:${hashed_message}`

<details>
<summary> Typescript example of how to construct a proper hmac authorization header value. </summary>

```ts
import { createHmac, createHash } from 'crypto'
/**
 * generateAuthorizationHeaderValue creates an AShirt-compatible authorization message for API communication
 * 
 * @param data.method The HTTP verb/method used in the request, in all caps (e.g. GET)
 * @param data.path The path part of the url, immediately following the hostname. Almost always starts with /api
 * @param data.date The current date/time, in GMT, and in RFC 1123 format
 * @param data.body The body to send. Not sending a body? Use an empty string instead
 * @param data.accessKey The access portion of your API key
 * @param data.secretKey The secret portion of your API key in no-encoding (raw bytes)
 * 
 * @returns a string in the format accessKey:hmacEncodedMessage
 */
function generateAuthorizationHeaderValue(data: {
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' // more methods with a similar naming style are possible
    path: string
    date: string // in RFC1123 format
    body: string
    accessKey: string
    secretKey: Buffer
}) {
  const stringBuff = Buffer.from(
    data.method + "\n" +
    data.path + "\n" +
    data.date + "\n"
  )
  // note that this isn't encoded -- the result is a series of raw bytes.
  const bodyDigest = createHash('sha256').update(data.body).digest()
  const message = Buffer.concat([stringBuff, bodyDigest])
  const hmacMessage = createHmac('sha256', data.secretKey)
    .update(message)
    .digest('base64')
  return `${data.accessKey}:${hmacMessage}`
}

// =============== Verify the output
const b64SecretKey =
  "DuvC7Wzpnsa2vtnOYw0RPGWeSdVB5L2L++PLpwGNb5yPQW47BoT5sohaMknU6Sh6a+0d/8dMh+wBEa2IPMMcNQ=="
const secretKey = Buffer.from(b64SecretKey, "base64")
const accessKey = "P4qRS5sa346iHWZBB53qzzNm"
const result = generateAuthorizationHeaderValue({
  method: "POST",
  path: "/api/operations",
  date: "Sun, 21 Oct 2018 12:16:24 GMT",
  body: '{"slug":"test-op","name":"Test Op"}',
  secretKey,
  accessKey,
})
console.log(result)
console.log("Does this match the expected value? ", result === 'P4qRS5sa346iHWZBB53qzzNm:RlbnBDbg5hj/foncSzOnfDWOCrTapyaL7fqKxkcCsFE=')
// expected output:
// P4qRS5sa346iHWZBB53qzzNm:RlbnBDbg5hj/foncSzOnfDWOCrTapyaL7fqKxkcCsFE=
// ================== Small helper to format the date in the right way
export function nowInRFC1123(): string {
  return new Date().toUTCString()
}
```

</details>

As of July 2022, the full API is supported via:

* [typescript](/enhancement_worker_templates/web/typescript_express/src/services/ashirt.ts)
* [javascript](/enhancement_worker_templates/lambda/js-container/app/ashirt_service.js)
* [python](/enhancement_worker_templates/web/python_flask/src/services/ashirt_base_class.py)

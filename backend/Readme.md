# AShirt Backend Server

A REST-based server for interfacing with the backend database and the AShirt frontend, screenshot client, and other API-based tools

## Overview

This project is a REST-based api server for the AShirt front end. The system is largely interacted with via `findings`, `evidence` and `operations`.

An :briefcase: _operation_ is the equivalent of a project, or root category. Each operation has a collection of evidence, and a collection of findings based on that evidence. Operations are scoped to a user but may be shared with other users.

:exclamation: _findings_ represents a particular vulerability or related point of interest that may need to be addressed by the operation team. It is comprised of evidence and a description. It also inherits tags from the associated evidence.

:mag: _Evidence_ is some media (e.g. images, code snippets) that helps build up some finding. Findings and Evidence live in a many-to-many relationship -- that is, findings can share evidence, and each finding is comprised of (potentially) multiple evidence

The server is broken into two services. The frontend communicates entirely with `web` routes, (see `server/web` for available routes), while automated services/tools communicate with `api` (see `server/api` for available routes)

### Configuration

Configuration is handled entirely via environment variables. To that end, here are the currently supported environment variables. Note: this section is likely to become out of date over time. Please check variables by searching the project for `os.Getenv` to discover all possible configurations.

* Environment Variables
  * `DB_URI`
    * MySQL connection string
  * `APP_IMGSTORE_BUCKET_NAME`
    * Assumes Amazon S3 storage location
  * `APP_IMGSTORE_REGION`
    * Amazon S3 region (e.g. us-west-1)
  * `APP_CSRF_AUTH_KEY`
    * The actual authorization key
    * Web Only
  * `APP_SESSION_STORE_KEY`
    * The actual session key
    * Web Only
  * `APP_PORT`
    * Configures what port the service starts on
    * Expected type: integer
  * `APP_RECOVERY_EXPIRY`
    * Specifies how long recovery codes/urls are active
    * Expected type: time duration (e.g. `60m` => 60 minutes `24h` => 24 hours)
    * Defaults to 24 hours
    * Base unit is 1 minute. Fractional minutes will be ignored
  * `APP_DISABLE_LOCAL_REGISTRATION`
    * Removes the registration aspect of local auth. Users will still be able to log in if they already have a local auth account.
    * Valid Options: `"true"` or `"false"`
    * Admins can provision new local auth accounts to provide access to new users.
    * Only valid if deploying using `ashirt` authentication. Otherwise has no effect
  * `AUTH_SERVICES`
    * Defines what authentication services are supported on the backend. This is limited by what the backend naturally supports.
    * Values must be comma separated (though commas are only needed when multiple values are used)
    * Example value: `ashirt,otka`
    * Currently valid values: `ashirt`, `okta`
      * This list will likely become outdated over time. Consult the authschemes directory for a better idea of what is supported.
  * `AUTH_${SERVICE}_` Variables
    * These environment variables are namespaced per Auth Service. Each of these is a specific field that can be used to pass configuration details to the authentication service. Note that `${SERVICE}` must be replaced with a proper string, expected in all caps. For example `AUTH_GITHUB`, `AUTH_ASHIRT`, `AUTH_GOOGLE`
    * `AUTH_${SERVICE}_CLIENT_ID`
      * For OAuth2 solutions. This provides a client ID value to the auth service
    * `AUTH_${SERVICE}_CLIENT_SECRET`
      * For OAuth2 solutions. This provides the corresponding secret
    * `AUTH_${SERVICE}_ISSUER`
      * For OAuth2 solutions. This essentially provides a URL to redirect the authentication process
    * `AUTH_${SERVICE}_BACKEND_URL`
      * The location of the ashirt service
    * `AUTH_${SERVICE}_SUCCESS_REDIRECT_URL`
      * For OAuth2 solutions. Where to redirect the user when login is successful
    * `AUTH_${SERVICE}_FAILURE_REDIRECT_URL_PREFIX`
      * Where to redirect the user when login fails for some reason. Note that this is a _prefix_. Current expected values are:
        * `/autherror/noverify`: User authentication failed (either challenge or token)
        * `/autherror/noaccess`: User authentication succeeded, but the user is excluded from using this application
        * `/autherror/incomplete`: User authentication succeeded and is able to use the application, but a matching ashirt user profile could not be created.

### Authentication and Authorization

Authentication is a somewhat modular system that allows for new authentication/identification to occur with external systems. The exact process is left pretty open to allow for maximum extensibility, while trying to keep a fairly simple interface. For details on how to add your own authentication scheme, see the [Custom Authentication](#custom-authentication).

Authorization is handled via the policy package. Policies are broken into two flavors: what operations can an authenticated user perform, and what operations can an authenticated user perform for a given operation. Each specific action is listed inside the policies, and each check happens prior to performing the requested action; generally, but not necessarily, these checks happen in the services package.

#### Administrator Priviledges

The AShirt backend and frontend have support for system administrator functions. Administrators gain priviledged access to some functionality, such as viewing and  deleting users, as well as managing operations. Administrators can bestow administrator status on any other user, and likewise can remove administrator access from any other user. This is all done, on the frontend, via an admin dashboard. On the backend, this is done via particular routes that verify admin status at the start of an admin-supported operation.

One limitation to this behavior is that, generally speaking, admins cannot alter themselves.

##### First Admin

When a fresh system is deployed, no users are present, thus no admins are present either. The first administration account, therefore, is granted to the first user that registers within the system.

#### Custom Authentication

Adding your own authentication is a 3 step process:

1. On the backend, create a new authscheme
   1. This is the bulk of the work. There are two interface methods to implement:
      1. `Name`: Every authentication needs a distinct name. The specific name does not really matter, but should be distinct from other utilized authentication scheme names.
         1. Note: Although the name does not matter, custom authentications **must not** use `,` in their names, as this is important for querying in some cases.
      2. `BindRoutes`: This provides a namespaced router that can be used to implement any routes needed to statisfy the authentication routine. In addition to the namespaced router a set of callback functions, called an AuthBridge, is provided to interact with the underlying system. Specifically, 3 functions have been provided to help provide access into the database: `CreateNewUser`, which attempts to instantiate a new _AShirt_ user into the database. `LoginUser`, which provides a mechanism for the backend to record a new session, and `FindUserAuthsByUserSlug`, which provides a mechanism to lookup existing users belonging to a specific identity provider (i.e. backing authscheme) and a user key (similar to a shortname or email, but specific to an authscheme).
2. The new authscheme needs to be "registered" so that the webserver will know to use it.
   1. Inside `bin/web.go`, create a new instance of the authscheme, then provide this as an argument to the `server.WebConfig` structure. Note that multiple authentication schemes can be present at once
3. The frontend needs to be updated to provide a way to login via your new authentication scheme, which is outside the scope of this miniguide.

#### Default AShirt authentication

Presently, at least some kind of authentication is required to use this service. AShirt provides a minimal authentication implementation to serve in this capacity. This implementation can be found in `authschemes/localauth/local_auth.go`

#### Account Recovery

Account recovery can be triggered by an admin for any user (except themselves). The account in question will generate a one-time-use code that expries in 24 hours. The user will need a special url that includes this code in order to login. Once logged in, the user will have full access to their account. At which point, they should probably link some other authentication system to their account, though this is not a requirement. The recovery scheme is baked into this system automatically, and cannot be disabled, except by recompiling the backend, and specifically removing the addition of this auth scheme.

### API Keys

As mentioned above, other services can iteract with the system, under the guise of some registered user, without requiring the user to login while using the tool. To do this, a user must first create an API key pair, and then associate these keys with the external tool (e.g. screenshot client).

## Development Overview

This project utilizes Golang 1.13 (with modules), interfaces with a MySQL database and leverages Gorilla Mux to help with routing. The project is testable via docker/docker-compose and is also deployed via docker.

### Development Environment

This project has been verified to build and run on Linux and MacOS X. Windows may work with some adjustments to supporting scripts. See the [dependencies](#dependencies) section for details on additional software for building. No specific IDE or editor is required, though there are some [notes](#visual-studio-code-notes) on integrating with [Visual Studio Code](https://code.visualstudio.com/)

### Dependencies

* Go 1.13
  * To get supporing libraries, use `go mod download`
  * To clean up libraries, use `go mod tidy`
* MySQL 8
  * This is started as part of the docker-compose script (meaning you won't actually need mysql locally), but all queries are targeted against this database system.
* Docker / Docker-compose
* Amazon S3 access (for production -- development versions use the `/tmp` directory)

### Buliding

Local binaries can be build via:

* api
  * `go build bin/api/*.go`
* web
  * `go build bin/web/*.go`

### Running the project

This project is best started in conjunction with the frontend and server. As such, a docker-compose file has been created to help launch all of the projects in the proper configuration. Inside the larger AShirt project is a `docker-compose.yml` file that can be started. Simply run `docker-compose up --build` to start this process.

Once the servers have been started, you can access the UI from `localhost:8080`. You can access the API from `localhost:3000`. The database lives on `localhost:3306`. Note that all end users (both from the website, and from tools utilizing the api) will interact with `localhost:8080/{service}`, with routing handled under the hood by external processes. By default, `localhost:8080/web` will direct the user to the web routes, while `localhost:8080/api` will direct the user to api routes. Any other routes will be interpreted by the frontend. No direct database access is provided to these users.

#### Notes

* The first run takes awhile to start, due to a number of required startup tasks. Subsequent runs should be quick.
* Changes to the database schema or switching branches _may_ require stopping the server and restarting it.
* The dockerfile is set up to hot reload changes, but given the way docker-compose restarts work, long periods spent debugging or making code changes may make the rebuild process take extra long. In these cases, it may be faster to stop and restart the docker-compose process manually.

### Using Seeded Data

Both unit tests and developer tests / manual tests use the same seed data to quickly spin up a decent selection of use cases within the database. This data is ever expanding, but in general tries to hit each of the features or expected bug scenarios. The most up-to-date document is going to be the seed data itself, which can be found at: `backend/database/seeding/hp_seed_data.go` (for a Harry Potter themed seed). However, a more pratical guide is as follows:

#### Using seed data for developer testing

* Several users are predefiend (see below). In general, the most "complete" users are:
  * Albus (Dumbledore) -- the super admin, indirect access to all operations
  * Ron (Weasley) -- admin for Chamber of Secrets
  * Harry (Potter) -- admin for Sorcerer's Stone
  * Draco (Malfoy) -- (mostly) no access, read-only access for Goblet of Fire
  * Nicholas (de Mimsy-Porpington) ; AKA: Nearly-Headless Nick -- A headless user. Note that Nick only has access to the Goblet of Fire operation
  * Tom (Riddle) -- deleted user
  * Rubeus (Hagrid) -- disabled user
* Users log in via their first name for their username and password. The password is always lowercase-only. e.g. Ron Weasley's login is `ron`/`ron`
* All users (except Tom Riddle) should see the Goblet of Fire operation
* The "Harry Potter and the Curse of Admin Oversight" operation provides no real evidence, but enough evidence to render a pattern on the operation overview page
* There is nuanced permission data for Sorcerer's Stone and Chamber of Secrets

#### Using seed data for unit testing

##### Setting up seeded data

Each test that wishes to use the seeded data needs to do the following:

```go
  db := seeding.InitTest(t) // this initializes the database connection to a fresh instance. This expects a certain path to the migrations directory, as well as a specific database name. See below for details on how to modify these
  err := seeding.HarryPotterSeedData.ApplyTo(db) // seeds the database with the harry potter seed data
  require.NoError(t, err) // ensure that no error was encountered while starting up
  userContext := seeding.SimpleFullContext(seeding.UserHarry) // Provide a proper authenticated policy for a given seed user. (note: any user can be used here -- Harry is just an example)

  // additional test-specific logic
```

This will spin up a fresh database instance the seeded data, and a user to perform the action (See users list below for pertinent details on seed users)

As a small caution, note that every time the database is refreshed, some time is spent establishing a new connection to the database and feeding the database both the schemea and a set of data. This process is relatively quick -- less than a second, but can quickly balloon once more tests are added.

##### Unit testing conventions

Unit tests should follow these guidelines:

* Ideal tests should verify access requirements for Read/Write, and Admin/Super Admin when necessary.
* Tests should use `testify.require` or `testify.assert` to validate condtions

#### Seeded Users

Note that this list may become out of date. Users with flags should be considered constant with respect to the below
fields, and Harry, Ronald, Hermione, Seamus, Ginny and Neville should be considered constant for the below fields as well.

| User                         | User key | Password   | Flags       | SS Permissions | CoS Permissions |
| ---------------------------- | -------- | ---------- | ----------- | -------------- | --------------- |
| Albus Dumbledore             | Albus    | `albus`    | Super Admin | Admin          | Admin           |
| Harry Potter                 | Harry    | `harry`    |             | Admin          | Write           |
| Ronald Weasley               | Ron      | `ron`      |             | Write          | Admin           |
| Hermione Granger             | Hermione | `hermione` |             | Read           | Write           |
| Seamus Finnegan              | Seamus   | `seamus`   |             | Write          | Read            |
| Ginny Weasley                | Ginny    | `ginny`    |             | NoAccess       | Write           |
| Neville Longbottom           | Neville  | `neville`  |             | Write          | NoAccess        |
| Draco Malfoy                 | Draco    | `draco`    |             | NoAccess       | NoAccess        |
| Serverus Snape               | Serverus | `serverus` |             | NoAccess       | NoAccess        |
| Cedric Digory                | Cedric   | `cedric`   |             | NoAccess       | NoAccess        |
| Fleur Delacour               | Fleur    | `fleur`    |             | NoAccess       | NoAccess        |
| Viktor Krum                  | Viktor   | `viktor`   |             | NoAccess       | NoAccess        |
| Alastor Moody                | Alastor  | `alastor`  |             | NoAccess       | NoAccess        |
| Minerva McGonagall           | Minerva  | `minerva`  |             | NoAccess       | NoAccess        |
| Lucius Malfoy                | Lucius   | `lucius`   |             | NoAccess       | NoAccess        |
| Sirius Black                 | Sirius   | `sirius`   |             | NoAccess       | NoAccess        |
| Peter Pettigrew              | Peter    | `peter`    |             | NoAccess       | NoAccess        |
| Parvati Patil                | Parvati  | `parvati`  |             | NoAccess       | NoAccess        |
| Padma Patil                  | Padma    | `padma`    |             | NoAccess       | NoAccess        |
| Cho Chang                    | Cho      | `cho`      |             | NoAccess       | NoAccess        |
| Rubeus Hagrid                | Rubeus   | `rubeus`   | Disabled    | NoAccess       | NoAccess        |
| Tom Riddle                   | Tom      | `tom`      | Deleted     | NoAccess       | NoAccess        |
| Nicholas de Mimsy-Porpington | Nicholas | `nicholas` | Headless    | NoAccess       | NoAccess        |

### Project Structure

The project contains various source code directories, effectively acting as a collection of mini-libraries interacting with each other.

```sh
├── authschemes                        # location for implemented authentication modules
│   ├── localauth                      # Location of authentication utilizing the base authentication system. Useful as an example if constructing custom authentication
│   └── {other auths as needed}        # recommended location for additional authentication schemes
├── bin                                # Main lines / build targets
│   ├── api                            # Target for building the api server
│   ├── dev                            # Code for _running_ the dev server
│   └── web                            # Target for building the webserver
├── config                             # Where server configuration details are parsed/how they're accessed
├── contentstore                       # Code providing abstraction over how to interact with remote media (specifically, images)
├── database                           # Code related to directly interacting with the database
├── dtos                               # Some DTOs. _Logical_ database structures (i.e. how you want to interact with the database)
├── helpers                            # A collection of pure functions used across different packages
├── integration                        # Integration tests
├── migrations                         # Contains all of the database changes needed to bring the original schema up to date
├── models                             # Exact("Physical") database structures (i.e. how you need to interfact with the database)
├── policy                             # _Authorization_ roles and rules to restrict access to APIs
├── server                             # Route endpoint definitions and basic request validation
│   ├── dissectors                     # A builder-pattern like solution for interpreting request objects
│   ├── middleware                     # Middleware to assist with request handling
│   ├── remux                          # A rewrapping package for better ergonmics when utilizing gorilla mux
│   ├── api.go                         # Routes for the "API" / screenshot tool
│   └── web.go                         # Routes for the web service
├── services                           # Underlying service logic. Also includes a number of unit tests
├── errors.go                          # Some helpers to build standard errors used across the system
├── Readme.md                          # This file!
├── run-dev.sh                         # Enables hot-relodaing of the dev server
└── schema.sql                         # The accumulated deployment schema -- useful when starting from scratch
```

### Errors and logging

The error model used within this application adopts the following principles:

* Use structured logging, to help finding/reporting errors
  * Logs are of the form: `timestamp=<ISO8601> key=value`
  * Common labels and meanings
    * `error` the error text for the underlying error. Wrapped errors are separated by ` : `
    * `msg` a general note on what operation is happening, or what unusual thing just happened
    * `ctx` the unique identifier that corresponds all (eligible) messages together by a particular request
    * All other values generally represent application state
* Use wrapped errors to help pinpoint the path an error took
* Export a (formatted) stacktrace for unexpected panics
* All error messages have two messages: a public one, exposing no real information to the user, and a private one, that gets logged
  * Errors containing the following text:
    * "Unwilling to" suggests that a request did not pass a permissions check.
    * "Unable to" suggests that some critical data was missing
    * "Cannot" suggests that we tried, and failed, to do the requested operation
    * messages that do not match the above generally have more specific information to identify them

### Visual Studio Code Notes

If you're using Visual Studio Code, you may want to make these changes:

1. Update your file associations for Dockerfile
   1. By default, the Docker plugin for vs code only provides a file association for `Dockerfile`. Since there are multiple dockerfiles here, if you want the files to be properly associated with the docker plugin, you should adjust your workspace or project configurations to include:

       ```json
       "files.associations": {
           "Dockerfile.*": "dockerfile"
       }
       ```

2. Recommended plugins:
   1. docker (ms-azuretools.vscode-docker)
   2. Go (ms-vscode.go)
3. Configuration settings:
   1. add this to your config to run all tests without error: `"go.testTimeout": "90s"`
      1. Running all tests can take some time. By default, VSCode's default timeout for running all tests is 30s. Since we have to reset the database between tests, our tests take a bit longer. 

### Common Tasks

* Updating the database schema
  1. Create a pair of migration files via `${PROJECT_ROOT}/bin/create-migration <name of change>`
  2. This will generate 2 files: a `up` version and a `down` version to reflect making the change and unmaking the change, respectively
  3. In the `up` version, provide the proper SQL statements to adjust the schema as needed
  4. In the `down` version, provide the opposite SQL statements to revert the changes
  5. While developing, **make sure that the database is running**, otherwise the next step will fail
  6. Once done with the pair of changes, run `${PROJECT_ROOT}/bin/migrate-up` to provide a new `${PROJECT_ROOT}/backend/schema.sql` file and update the running database
  
  Note: you may also need to update the `models` and/or the `dtos`

## Contributing

TBD

## License

TBD

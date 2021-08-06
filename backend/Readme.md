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
  * `APP_FRONTEND_INDEX_URL`
    * Used by the backend to redirect to the frontend in some scenarios (e.g. Email-based recovery)
  * `APP_BACKEND_URL`
    * Common field used for some authentication schemes. Provides a mechanism for the backend to reference itself to authentication providers
  * `APP_SUCCESS_REDIRECT_URL`
    * Used in some authentication schemes to redirect to the frontend after a successful authentication.
  * `APP_FAILURE_REDIRECT_URL_PREFIX`
    * Used in some authentication schemes to redirect to the frontend after a failed authentication. 
  * `AUTH_SERVICES`
    * Defines what authentication services are supported on the backend. This is limited by what the backend naturally supports.
    * Values must be comma separated (though commas are only needed when multiple values are used)
    * Example value: `ashirt,otka`
    * Currently valid values: `ashirt`, `okta`
      * This list will likely become outdated over time. Consult the authschemes directory for a better idea of what is supported.
  * `AUTH_${SERVICE}_` Variables
    * These environment variables are namespaced per Auth Service. Each of these is a specific field that can be used to pass configuration details to the authentication service. Note that `${SERVICE}` must be replaced with a proper string, expected in all caps. For example `AUTH_GITHUB`, `AUTH_ASHIRT`, `AUTH_GOOGLE`
    * `AUTH_${SERVICE}_CLIENT_ID`
      * This provides a client ID value to the auth service
      * For OIDC and Okta authentication
    * `AUTH_${SERVICE}_CLIENT_SECRET`
      * This provides the corresponding secret
      * For OIDC and Okta authentication
    * `AUTH_${SERVICE}_ISSUER`
      * This essentially provides a URL to redirect the authentication process
      * For Okta authentication
      * Deprecated
    * `AUTH_${SERVICE}_BACKEND_URL`
      * The location of the ashirt service
      * For Okta Authentication
      * Deprecated
    * `AUTH_${SERVICE}_SUCCESS_REDIRECT_URL`
      * Where to redirect the user when login is successful
      * For Okta Authentication
      * Deprecated
    * `AUTH_${SERVICE}_FAILURE_REDIRECT_URL_PREFIX`
      * Where to redirect the user when login fails for some reason. Note that this is a _prefix_. Current expected values are:
        * `/autherror/noverify`: User authentication failed (either challenge or token)
        * `/autherror/noaccess`: User authentication succeeded, but the user is excluded from using this application
        * `/autherror/incomplete`: User authentication succeeded and is able to use the application, but a matching ashirt user profile could not be created.
      * For Okta Authentication
      * Deprecated
    * `AUTH_${SERVICE}_TYPE`
      * Supported Values: `oidc` (Note that `local` and `okta` are reserved values, and not usable)
      * Required for all authentication types
    * `AUTH_${SERVICE}_NAME`
      * Must be distinct among auth service names
      * For OIDC authentication
    * `AUTH_${SERVICE}_FRIENDLY_NAME`
      * The name of the authentication scheme presented to the end user
      * For OIDC authentication
    * `AUTH_${SERVICE}_SCOPES`
      * Used to help pull additional scopes, which would be useful if the standard scopes are insufficient.
      * At a minimum, the `openid` and `profile` scopes are requested.
      * For OIDC authentication
    * `AUTH_${SERVICE}_PROVIDER_URL`
      * Used to help point to the OIDC provider's discovery document. Note that this URL _MUST_ match the issuer value in the discovery document.
      * For OIDC authentication
    * `AUTH_${SERVICE}_PROFILE_FIRST_NAME_FIELD`
      * Used within the application to refer to the user's first name. This is only used as an intitial value. Can be updated in the user's settings
      * Optional. Defaults to `given_name` (a common claim type)
      * For OIDC authentication
    * `AUTH_${SERVICE}_PROFILE_LAST_NAME_FIELD`
      * Used within the application to refer to the user's last name. This is only used as an intitial value. Can be updated in the user's settings
      * Optional. Defaults to `family_name` (a common claim type)
      * For OIDC authentication
    * `AUTH_${SERVICE}_PROFILE_EMAIL_FIELD`
      * This is used to as a mechanism to contact the user via email (currently only used for recovery)
      * Optional. Defaults to `email` (a common claim type)
      * For OIDC authentication
    * `AUTH_${SERVICE}_PROFILE_SLUG_FIELD`
      * This is functionally equivalent to a username or an email for most services. Used internally for associating a user to their content and assignments
      * Must provide a unique value for all users using this authentication scheme.
      * Optional. Defaults to `email` (a common claim type)
      * For OIDC authentication
    * `AUTH_${SERVICE}_DISABLE_REGISTRATION`
      * Prevents new registrations for the given service
      * Optional. Defaults to `false` (open registration)
      * For OIDC authentication
  * `EMAIL_FROM_ADDRESS`
    * The email address to use when sending emails. The specific value may be influenced by your email provider
  * `EMAIL_TYPE`
    * Indicates what kind of email service is used to send the emails.
    * Valid values: `smtp`, `memory` (for test), `stdout` (for test)
  * `EMAIL_HOST`
    * The location of the email server. If connecting to an SMTP server, a port is also required (e.g. `my-email-server:25`)
  * `EMAIL_USER_NAME`
    * The username to use when authenticating with PLAIN or LOGIN SMTP servers
  * `EMAIL_PASSWORD`
    * The password to use when authenticating with PLAIN or LOGIN SMTP servers
  * `EMAIL_IDENTITY`
    * The identity to use when authenticating with PLAIN SMTP servers
  * `EMAIL_SECRET`
    * The secret to use when authenticating with CRAM-MD5 SMTP servers
  * `EMAIL_SMTP_AUTH_TYPE`
    * Indicates which kind of authentication scheme to use when connecting to an SMTP server
    * Valid values: `login`, `plain`, `crammd5` (for LOGIN, PLAIN, and CRAM-MD5 respectively)

### Authentication and Authorization

Authentication is a somewhat modular system that allows for new authentication/identification to occur with external systems. The exact process is left pretty open to allow for maximum extensibility, while trying to keep a fairly simple interface. For details on how to add your own authentication scheme, see the [Custom Authentication](#custom-authentication).

Authorization is handled via the policy package. Policies are broken into two flavors: what operations can an authenticated user perform, and what operations can an authenticated user perform for a given operation. Each specific action is listed inside the policies, and each check happens prior to performing the requested action; generally, but not necessarily, these checks happen in the services package.

#### Administrator Priviledges

The AShirt backend and frontend have support for system administrator functions. Administrators gain priviledged access to some functionality, such as viewing and  deleting users, as well as managing operations. Administrators can bestow administrator status on any other user, and likewise can remove administrator access from any other user. This is all done, on the frontend, via an admin dashboard. On the backend, this is done via particular routes that verify admin status at the start of an admin-supported operation.

One limitation to this behavior is that, generally speaking, admins cannot alter themselves.

##### First Admin

When a fresh system is deployed, no users are present, thus no admins are present either. The first administration account, therefore, is granted to the first user that registers within the system.

###### First Admin alternative

In certain situations, there may not be a way for a new user to register with AShirt without an
admin's help, even for the first user. In these cases, the below SQL can be used to create an initial
account and a recovery code to link the account to a supported authentication scheme.

Note that this requires direct access to the database. This should only be done for the first user
when the normal approach will not work.

1. Edit, and execute the below SQL

  ```sql
  INSERT INTO users (slug, first_name, last_name, email, admin) VALUES
  ('user@example.com', 'User', 'McUserface', 'user@example.com', true);

  INSERT INTO auth_scheme_data (auth_scheme, user_key, user_id) VALUES
  ('recovery', 'e3c6ead16e0c25820ba730f278ef54133da5610f9bf1d2e481ff6693c8df85123a29b8dc1f033a2f', 1);
  ```

  This will add a one-time password to AShirt which will allow the admin to sign in. Note that,
  per convention, the slug and email should match if using ASHIRT Local Authentication. This is not
  a hard requirement if you want to deviate from the convention. All other fields can be updated by
  updating the profile in Account Settings.

2. Start up the AShirt frontend and backend, if not already started
3. Once started, edit, and navigate to: `http://MY_ASHIRT_DOMAIN/web/auth/recovery/login?code=e3c6ead16e0c25820ba730f278ef54133da5610f9bf1d2e481ff6693c8df85123a29b8dc1f033a2f`

The admin should now be logged in, and can update their security information.

1. Click the person icon and select "Account Settings"
2. Go to "Authentication Methods"
3. Find a supported login the admin wishes to use, and click the "Link" button. Follow this process.
   1. Note: if linking to ASHIRT Local Authentication, when the admin logs in, they will log in via the email address provided during the linking step, not (necessarily) the above sql script.

At this point, a proper admin account exists and you can log in via the linked methods.

#### Open ID Connect (OIDC) Authentication

Authentication via OIDC is supported under the condition that the ODIC provider have a discovery document. A discovery document provides the urls necessary for the implementation to interact autonomously with the ODIC provider. An example of a discovery document can be found [here](https://accounts.google.com/.well-known/openid-configuration)

##### Adding an OIDC authentication provider

Each OIDC provider follows the same process:

1. In the `AUTH_SERVICES` environment variable, provide a new short name for the service. The name choice here is arbitary, but should be a single word (with underscores, if desired). The case used here is irrelevant. For our example, we will choose `pro_auth` as our key
2. Each OIDC authentication will need a number of environment variables with specific names to complete the configuration. The environment variables meaning is detailed [here](#configuration), but briefly, each key must be prefixed with `AUTH_${SERVICE}`, and it's meaning will be detailed below. In our case, since our service name is `pro_auth`, our prefix will be `AUTH_PRO_AUTH` and the expected values are:

  ```sh
    AUTH_PRO_AUTH_TYPE: oidc                                # Flags to the backend that OIDC authentication should be used
    AUTH_PRO_AUTH_NAME: pro_auth                            # The name of the service within the database. Can be anything, but it's recommended that it be the same as the auth_service value.
    AUTH_PRO_AUTH_FRIENDLY_NAME: ProAuth                    # The name of the service, as presented to the user (e.g. in this case, they'll see a button with the text "Login with ProAuth")
    AUTH_PRO_AUTH_CLIENT_ID: clientID123                    # The client ID provided by the OIDC provider.
    AUTH_PRO_AUTH_CLIENT_SECRET: sup3rs3cr3tK3y             # The client secret provided by the ODIC provider.
    AUTH_PRO_AUTH_SCOPES: email                             # What additional scopes to load when getting an identity token. For most services, this can be "email". 
    AUTH_PRO_AUTH_PROVIDER_URL: https://myacct.proauth.com  # The provider URL for your service. In general, this should be the "issuer" field specified in the discovery document. Convieniently, you can also test this value by adding "/.well-known/openid-configuration" to the end of the URL and seeing if the concatinated value produces a discovery document. If so, then this is likely the provider url
  ```

3. In most cases, the above should be sufficient to have a working OIDC implementation. However, it may be necessary in some instances to provide some additional configuration. This is because after getting a new login, we need to create a user account for AShirt, which requires some personal info -- specifically, a first and last name, email, and another unique value (which can also be email, if desired). You can use the below fields to customize/complete your experience.

  ```sh
  AUTH_PRO_AUTH_PROFILE_FIRST_NAME_FIELD: first_name  # Retrieve the "first name" value from the named claim
  AUTH_PRO_AUTH_PROFILE_LAST_NAME_FIELD: last_name    # Retrieve the "last name" value from the named claim
  AUTH_PRO_AUTH_PROFILE_EMAIL_FIELD: email            # Retrieve the "email" value from the named claim
  AUTH_PRO_AUTH_PROFILE_SLUG_FIELD: username          # Retrieve the "slug" value from the named claim -- used to uniquely identify a user within the system -- note that typically, email is sufficient, but other options may be available in your identity provider.
  ```

4. Finally, OIDC authentication supports registration lockouts. In this scenario, registration will be denied for all new users that do not currently have a login using that authentication scheme. This does not prevent users from linking that authentication type, only preventing completely new accounts. This will be most useful for public OIDC providers (e.g. Google oidc) that cannot limit access via a user's list. To disable registration, using our example configuratoin, we would accomplish this via: `AUTH_PRO_AUTH_DISABLE_REGISTRATION: "true"`. If registration is disabled, you can still invite users via a small workaround. See [here](#recovery-based-user-invites-workaround) for the workaround details.

##### Provider URLs

Here are some provider urls for some common OIDC providers

| Service  | URL                                               |
| -------- | ------------------------------------------------- |
| Okta     | https://${Okta-client-ID}.okta.com                |
| Google   | https://accounts.google.com                       |
| OneLogin | https://${Onelogin-client-ID}.onelogin.com/oidc/2 |

##### Migrating from Okta to generic OIDC Okta

The original Okta authentication instance has changed. Okta is still supported but the custom
integration is now deprecated and it is now recommended that Okta integration is accomplished by
using generic OIDC. Here's a mini-guide on performing that conversion.

This guide assume that your okta authentication (located in `AUTH_SERVICES` is called "okta". If it is not "okta" then each of the environment variables will be slightly different. For example, if your okta instance is called "my_okta" then your "AUTH_OKTA_TYPE" would actually be called "AUTH_MY_OKTA_TYPE"

1. Start with the base configuration:

   ```sh
   AUTH_OKTA_TYPE: oidc           # Specifies that this uses OIDC authentication
   AUTH_OKTA_NAME: okta           # This is a name internal to the application -- must be unique
   AUTH_OKTA_FRIENDLY_NAME: Okta  # This is the name presented to the user when they see the login button
   AUTH_OKTA_SCOPES: email        # This specifies to load the "email" scope in addition to the standard scopes
   ```

2. The `AUTH_OKTA_CLIENT_ID` and `AUTH_OKTA_CLIENT_SECRET` fields are unchanged, and can simply be left alone.
3. Create the `AUTH_OKTA_PROVIDER_URL` with the value from `AUTH_OKTA_ISSUER`. This value need to be updated. Simply remove the `/oauth2/default` portion of the Issuer URL to create the provider URL. For example, given the issuer URL `https://MY_OKTA_INSTANCE.okta.com/oauth2/default`, the provider value will be `https://MY_OKTA_INSTANCE.okta.com`. 
4. The following fields move from Okta-specific configurations to common configurations. Simply rename the environment variable as follows:
   * `AUTH_OKTA_BACKEND_URL` => `APP_BACKEND_URL`
   * `AUTH_OKTA_SUCCESS_REDIRECT_URL` => `APP_SUCCESS_REDIRECT_URL`
   * `AUTH_OKTA_FAILURE_REDIRECT_URL_PREFIX` => `APP_FAILURE_REDIRECT_URL_PREFIX`
5. Finally, the `AUTH_OKTA_PROFILE_TO_SHORTNAME_FIELD` has been renamed to `AUTH_OIDC_OKTA_PROFILE_SLUG_FIELD`. Simply rename the field and keep the existing value.

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

A separate set of recovery exists for users to initiate a self-service recovery. In this case, users will need to select the "Forgot your password?" option from the login page. This method is expected to only be valid for local/default loigin. Users will receive an email with a link to recover their account. The recover code will expire in 24 hours from the time the email was sent.

#### Preprovisioning / Inviting users

In certain circumstances, you may want to create an account for a user you anticipate joining. Admins can do this via navigating to "User Management" on the frontend admin console, and clicking the "Create new user" button. This will create a new local account, and provide the admin with a one-time login for the new user.

##### Recovery-based user invites (Workaround)

In certain situations, having a local auth user account may not be ideal, but you may still want to preprovision a new user. This is possible via a small workaround with some existing functionality. See the below for the steps.

Note: Local Authentication must still be enabled in this situation, even if it is not used.

1. Login as an admin
2. Navigate to the admin tools, and specifically to the User Management screen
3. Click on the `Create New User` button to create an initial user account. Provide valid data for the existing fields, and remember the name given
4. After creating the new user, search for that user in the User List.
5. Once you find the user, under `Actions`, choose the `Generate Recovery Code`
6. Provide the recovery URL to the new user. they can use this to do a one-time login. Along with the code, tell the new user to link their account via one of the approved authentication methods.
  
After this, the user will be able to login normally, using their preferred login mechanism.

Note that the one-time login via local auth will still be active.

To remove the one-time password:

1. Find the user in the User List
2. Choose `Edit User`, and navigate to `Authentication Methods`
3. Find the `local` authentication scheme, and under Actions, choose `Delete`

### API Keys

As mentioned above, other services can iteract with the system, under the guise of some registered user, without requiring the user to login while using the tool. To do this, a user must first create an API key pair, and then associate these keys with the external tool (e.g. screenshot client).

### Emails

The backend has a system to send emails out to notify users (with an email address) as needed. Currently, this system is only used to send account recovery emails. An email server will be needed, but stmp services can be configured via environment variables.

Custom email services can be implemented or extended by meeting the `EmailServicer` interface in `emailservices/interface.go`. 

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

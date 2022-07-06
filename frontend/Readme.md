# AShirt Frontend

A SPA website to allow users and teams to track security issues and vulnerabilities via storing evidence, combining evidence into findings, and organizing based on operations.

## Overview

This project exists as a frontend to the AShirt backend/services component of the larger AShirt project. In general, users login, can create _operations_, add _evidence_ to an operation, invite other users to contribute to that operation, organize evidence into larger _findings_, and then share those findings with others. All of the above is governed by some simple permission rules, allowing users to control who and how users can contribute or view an operation.

:briefcase: Operations are user-owned projects. The user that creates an operation is the admin for that operation, and can set the name, status, tags, and user associated with that operation. The owner of a project can grant users 3 levels of access: `read`, where users are able to see the content of an operation. `write`, where users a can add/remove/edit evidence in an operation (and also implicitly `read` the operation), and `admin`, which can make further changes to the operation settings (and also implicitly `read` and `write` to the project).

:mag: Evidence is some media (E.g. code snippets, images, accounts) that is added to an operation. Currently, all operations accept all kinds of evidence. Evidence is dated and optionally tagged with a description of the event. Users may add new tags as needed, while operation administrators can alter or remove added tags as necessary.

:exclamation: Findings are grouped sets of evidence with a title, category and optionally a description of the finding. Current finding categories are: [Product, Network, Enterprise, Vendor, Behavioral, Detection Gap]. Descriptions support Markdown. Findings also contain tags, but tags are aggregated from associated evidence, rather than having tags specifically for findings.

In addition to the above, the following features are supported:

* Admin accounts
  * Admins have special access to the system. They can view all users, reset password, and edit any operation settings.
  * Admins gain access to these functions, and more, by choosing the person icon, and selecting "Admin"
* Robust searching
  * Sort by tag, date range or description
* Saving searches for evidence and findings (on a per-operation basis)
* Lightbox view of all evidence, with hotkeys for navigating between evidence without leaving the lightbox
  * Left Arrow / Up Arrow / J to navigate to more recent content
  * Right Arrow / Down Arrow / K to navigate to older content
* Generating personal API keys for related tooling

## Deployment flags

This service can change its rendering based on the flags provided _to the backend_. See the backend details on flags [here](/backend/Readme.md#flags).

## Development Overview

This project utilizes Typescript 3.5 and React 16.8 to construct a versitle and robust website. Packaging is handled via webpack, while dependencies are handled via npm. No specific IDE is required to develop this application.

### Dependencies

* Node 12+
  * `npm`
* Docker / Docker-compose (soft requirement. The frontend can be started locally as needed, though pairs best with an entire system)
* Various npm dependencies, detailed in `package.json`
* [Stylus](http://stylus-lang.com)

Dependencies can be retrived via `npm install`.

### Building

Local builds can be built via `npm build`, as usual.

### Running the project

To run the project on the host machine, you can run `npm start`, as usual.

To run an entire system at once, utilize docker-compose file located in the larger AShirt directory via `docker-compose up --build`. The `--build` flag is important here when dependencies are added/removed. Once the system intializes and the server is available, the frontend can be found on `localhost:8080`

### Project Structure

```sh
├── public                            # Home of the index.html file
├── src                               # Where all of the typescript code is stored
│   ├── base_css                      # Starting point for css
│   ├── components                    # Custom components are stored here
│   ├── helpers                       # Pure functions that are used throughout the codebase
│   ├── pages                         # Screens and modals within the SPA
│   │   ├── account_settings          # Handles users settings
│   │   ├── admin                     # Handles admin settings
│   │   ├── admin_modals              # Modals presented to admins (or, more specifically, generated on admin pages)
│   │   ├── login                     # The login page
│   │   ├── not_found                 # The 404 page
│   │   ├── operation_edit            # The settings page for operations
│   │   ├── operation_list            # The page where users can view operations visible to them / create new operations
│   │   └── operation_show            # The page where users see the details of an operation (evidence, findings)
│   ├── services                      # Provides an interface into external services. In other words, provides the logic to make async requests to the backend.
│   ├── auth_context.tsx              # Provies a react context for storing the currently logged in user and their CSRF token
│   ├── global_types.ts               # Common, custom types used throughout the application
│   ├── index.tsx                     # React initialization file / root component
│   ├── routes.tsx                    # React-Router router for _most_ routes. Some subpages utilize a different scheme.
│   └── vars.styl                     # Common variables used in Stylus files
├── Dockerfile                        # A production-oriented docker file
├── Dockerfile.dev                    # A development-oriented docker file
├── nginx.conf                        # 
├── package.json                      # 
├── package-lock.json                 # 
├── Readme.md                         # This file!
├── tsconfig.json                     # Typescript convention/configuration file.
└── webpack.config.js                 # 
```

#### Component structure

Components have a pretty consistent structure. For example:

```sh
├── checkbox               # The name of the component
│   ├── check.svg          # Images, if any, that the component may use
│   ├── index.tsx          # The component display and typical location for logic (though sometimes more complicated logic is broken out into a separate file)
│   └── stylesheet.styl    # An overall stylesheet for the component. This is run through Stylus to convert to CSS.
```

more complicated components typically break out parts ito separate files or subdirectories, depending on if this is a specialized subcomponent/page, or for complicated logic, etc.

### Adding (Additional) Authentication

While the backend largely controls what authentication methods are supported, the frontend has a part in this as well. In general, authentication methods are handled in `src/pages/login/index.tsx`, with localauth login already provided. Additional buttons/components can be added here to support additional services.

### Notes

1. Classnames
   1. During production, classnames get mangled. This is a side effect of running these through stylus. In general, this should not be a problem for production
   2. During development, classnames are set to the path of the rendering component, to help distguish where some styles are introduced. If you look in the development console, you should see classnames like: `src/components/card/stylesheet__root src/pages/operation_list/operation_card/stylesheet__root`. This indicates that the styles are picked up from `src/components/card/stylesheet` (specifically the root class), `src/pages/operation_list/operation_card` (specifically the root class)
      1. If you do not see a classname like the above, this is likely due to a classname that does not have a definition in stylesheet.styl
2. Some images used here have been pulled from icofont or fontawesome
   1. TODO: these need proper attribution

## Contributing

TBD

## License

TBD

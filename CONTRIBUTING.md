# Prerequisites
Before we can start, there are a few tools we need:

- [Go 1.16+](https://go.dev)
- [NodeJS 14+](https://nodejs.dev)
- [Docker](https://www.docker.com/get-started)
- [Make](https://www.gnu.org/software/make/) _(Optional)_

# Configuration
Before Maestro can run in, there are a few configuration files we need to set up.

## API
In the `configs/` directory, there is a `maestro.example.yaml` file, copy this to `maestro.yaml`.
From there, you can edit the configuration as required.

### Acquiring API keys
🚧 Todo 🚧

## Frontend
The frontend (located in `web/frontend-remix`) also has an `example.env` file which just contains the API url.
Copy this to `.env` and edit it as required. (Note that these are overridden in docker compose)

## Docker compose
You may have noticed an `example.env` file in the `configs/` directory. This can be used by docker compose to inject
environment variables. It's primarily used for storing secrets such as streaming service API keys and database credentials.

If you're looking to run the API and/or frontend via docker compose, it's recommended to copy the `example.env` file to
`.env`, and fill in the blanks as required.

Note that the `MAESTRO_` environment variables are not required, and these secrets _can_ be configured in the
`maestro.yaml` file if desired.

## Database
The `deployments/docker-compose.yaml` file provides a MongoDB container out of the box.
Provided that the `.env` file has been filled out correctly, this should work out of the box.
It's recommended to use this for development.

# Makefile
A `Makefile` is available in the root directory with a handful of useful commands, such as `make compose-fresh`, which
will automatically run `docker-compose` with the appropriate flags and arguments.

Run `make` or `make help` for a list of available commands and what they do.

# Kubernetes
Kubernetes is used to deploy Maestro to production. For a detailed guide on configuring Kubernetes and Maestro, head
over to the [Kubernetes readme](deployments/k8s/README.md).

# Pull Requests
If you have some changes you'd like to see merged into Maestro, consider forking and submitting a pull request!

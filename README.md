<h1 align="center">
  🎵 Maestro 🎵 
</h1>

<h3 align="center">
  A small web application and API that allows music to be shared across streaming services.
</h3>

[//]: # (Todo: Preview image)

# Configuration
Before Maestro can run, there are a few things which need to be configured.

## API
In the `./configs/` directory, there is a `maestro.example.yaml` file, copy this to `maestro.yaml`.
From there, you can edit the configuration as required.

### Docker compose
You may also notice an `example.env` file in the `./configs/` directory. This can be used by docker compose to inject
environment variables. It's primarily used for storing secrets such as streaming service API keys and database credentials. 

If you're looking to run the API via docker compose, it's recommended to copy the `example.env` file to `.env`, and fill
in the blanks as required.

Note that this is not required, and these secrets _can_ be configured in the `maestro.yaml` file if desired.

## Frontend
The frontend also has an `example.env` file which just contains the API url. Copy this to `.env` and edit it as required.

Note that this is overridden in docker compose.

# Contributing
Contributions are what make the open source community such an amazing place to be, learn, inspire, and create.
Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`feature/AmazingFeature`)
3. Commit your Changes
4. Push to the Branch
5. Open a Pull Request

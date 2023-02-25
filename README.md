<h1 align="center">
  🎵 Maestro 🎵 
</h1>

<h3 align="center">
  A lightweight web application for sharing music across a variety of music streaming platforms.

  [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/yukitsune/maestro/CI.yml?branch=main)](https://github.com/YuKitsune/Maestro/actions/workflows/CI.yml)
  [![Go Report Card](https://goreportcard.com/badge/github.com/yukitsune/maestro)](https://goreportcard.com/report/github.com/yukitsune/maestro)
  [![License](https://img.shields.io/github/license/YuKitsune/Maestro)](https://github.com/YuKitsune/Maestro/blob/main/LICENSE)
  [![Latest Release](https://img.shields.io/github/v/release/YuKitsune/Maestro?include_prereleases)](https://github.com/YuKitsune/Maestro/releases)

  <img src="Screenshot.png" />
</h3>

[//]: # (Todo: Link to showcase)
> **Warning**
> Maestro is not under active development. I'll still keep an eye out for any issues or PRs, but no new features or streaming services will be added.
> Maestro was originally created to encourage myself to learn new technologies.
> It's time to put this to bed and move onto the next side-project.

# What can it do?
Maestro aggregates links to artists, albums, and tracks across a number of different streaming services.
This lets you share music with anyone regardless of their preferred music streaming service (as long as it's either Spotify, Apple Music, or Deezer 😅)

# Configuration
In the `configs/` directory, there is a `maestro.example.yaml` file, copy this to `maestro.yaml`.
From there, you can edit the configuration as required.

The frontend (located in `web/frontend-remix`) also has an `example.env` file which just contains the API url.
Copy this to `.env` and edit it as required. (Note that these are overridden when running in docker compose)

## Docker compose
You may have noticed an `example.env` file in the `configs/` directory. This can be used by docker compose to inject
environment variables. It's primarily used for storing secrets such as streaming service API keys.

If you're looking to run the API and/or frontend via docker compose, copy the `example.env` file to `.env`, and fill in
the blanks.

Note that the `MAESTRO_` environment variables are not required, and these secrets _can_ be configured in the
`maestro.yaml` file if desired.

## Database
The `docker-compose.yaml` file provides a MongoDB container out of the box.
Provided that the `.env` file has been filled out correctly, this should work out of the box.

## Acquiring API keys

### Apple Music
Apple Music have a [guide](https://developer.apple.com/documentation/applemusicapi/getting_keys_and_creating_tokens) on acquiring the required keys.
Once you have the keys, you can use a tool like [amjwt](https://github.com/YuKitsune/amjwt) to generate the token (disclaimer: I wrote it).
Once you've created the token, copy it into config files mentioned above.

### Spotify
You'll need to create a new application using your Spotify account. You can visit [this page](https://developer.spotify.com/dashboard/applications) to get started.
Once you've created the application, make sure you copy the Client ID and Client Secret into config files mentioned above.

### Deezer
Deezer doesn't require any API keys

### Keeping your keys safe
As long as you keep your keys in the `maestro.yaml` and/or `.env` files, or even somewhere outside the repository, they
should be relatively safe.
Don't add them to the example config files, or any other checked in files. Make sure you review your changes before
accidentally committing your keys.

# Contributing
If you have some changes you'd like to see merged into Maestro, consider forking and submitting a pull request!

# Support
If you want to support this tiny sandbox project of mine, feel free to [buy me a coffee](https://www.buymeacoffee.com/yukitsune256)!
If there is enough interest, I may consider dedicating more time to this.

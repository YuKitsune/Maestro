# Maestro Helm Chart
Maestro is a small web application and API that allows lets you share the music you love across a variety streaming services.

# Installing
## Prerequisites
Before we can start, there are a few tools we need:

- [Helm 3+](https://helm.sh/docs/intro)
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

Add Maestro's chart repository to Helm
```shell
TODO
```

You can update the chart repository by running:
```shell
helm repo update
```

## Deploying
```shell
TODO
helm install maestro yukitsune/yukitsune --values ./values.yaml
```

### Specifying values
Maestro can be configured using the values file.

Example:
```yaml
# The publicly accessible URL for the Maestro API
# Example: "https://maestro.yukitsune.dev/api"
publicApiUrl: 

database:
  mongoUri: # The URI for the MongoDB database
  mongoDatabaseName: # The name of the MongoDB databse to use
  mongoTlsCert: # (Optional) The MongoDB server TLS certificate

services:
  appleMusic:
    token: # The Apple Music JWT

  spotify:
    clientId: # The Spotify client ID
    clientSecret: # The Spotify client secret
```

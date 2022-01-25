# Kubernetes Deployment
Before these Kuberneted configuration files can be used, a few changes need to be made.

## `secrets/`
The `secrets/` directory has some "blank slates" ending in `example.yaml`.
Simply copy these files, removing the `example` part from the file name, e.g:
```shell
cp maestro-api.secrets.example.yaml maestro-api.secrets.yaml 
```

From there, open them in your editor of choice and fill in the blanks.

### Traefik args
The `traefik-args.secrets.example.yaml` is slightly different, as it has templates for other Traefik certificates resolvers.
Copy the file as above, then uncomment and configure it to your hearts content!

## Hosts

### `host.config.yaml`
The `host.config.yaml` file will need to be updated with your own domain, as it's used in the frontend deployment.

_Ideally_ this would be the only file you need to update, but Kubernetes doesn't appear to support string interpolation,
nor does it allow for **any** value to be substituted from a ConfigMap. 😔

### `maestro-api` and `maestro-frontend` yaml files
At the bottom of the `maestro-api.yaml` and `maestro-frontend.yaml` files is the Traefik ingress config.
You'll need to update `spec.routes[0].match` to use your own domain.

## Setting up the cluster

Apply the ConfigMaps and Secrets
```shell
kubectl apply -f secrets/cloudflare.secrets.yaml -f secrets/maestro-api.secrets.yaml
kubectl apply -f host.config.yaml -f maestro-api.config.yaml
```

Install Traefik using Helm
```shell
helm install traefik traefik/traefik -f secrets/traefik-args.secrets.yaml -f traefik-values.yaml
```

Apply the Maestro services
```shell
kubectl apply -f maestro-api.yaml -f maestro-frontend.yaml
```

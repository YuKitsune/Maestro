# Kubernetes Deployment
In production, Maestro is deployed to a Kubernetes cluster. It's overkill, but it's also cool. 😎

## Prerequisites
The only things you need to deploy Maestro are [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) and [`helm`](https://helm.sh).

If you want to run a local Kubernetes cluster, you can use [`minikube`](https://minikube.sigs.k8s.io/docs/).

## Configuration
Before these Kubernetes configuration files can be used, a few config changes need to be made.

### Secrets
The `secrets/` directory has some "blank slates" ending in `example.yaml`.
Simply copy these files, removing the `example` part from the file name, e.g:
```shell
cp maestro-api.secrets.example.yaml maestro-api.secrets.yaml 
```

From there, open them in your editor of choice and fill in the blanks.

#### Traefik args
The `traefik-args.secrets.example.yaml` is slightly different, as it has templates for other Traefik certificates resolvers.
Copy the file as above, then uncomment and configure it to your hearts content!

These are kept separate as there can be some sensitive data in the additional arguments.

### Hosts

#### `host.config.yaml`
The `host.config.yaml` file will need to be updated with your own domain, as it's used in the frontend deployment.

_Ideally_ this would be the only file you need to update, but Kubernetes doesn't appear to support string interpolation,
nor does it allow for **any** value to be substituted from a ConfigMap. 😔

### `maestro-api` and `maestro-frontend` yaml files
At the bottom of the `maestro-api.yaml` and `maestro-frontend.yaml` files are the Traefik ingress configs.
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

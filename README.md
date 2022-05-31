# What is it

Inject secrets from Azure Vault to a process running inside a container in a kubernetes pod.

## How it works

1. mount the volume from which to read a 'template file'. this will be parsed and used to create an output file in which the secrets will be stored
1. mount an empty volume (this will be shared with the microservice container that will load the secrets)
1. connect to an azure keyvault specified in the `VAULT_URI` env var
1. get all the secrets, take the `INPUT_TEMPLATE_FILE` file template, parse and save the artifact in the file listed in `OUTPUT_FILE_PATH`
1. init container dies, leaves the `OUTPUT_FILE_PATH` file in the 'previously empty volume'
1. when the microservice starts, it will find the `OUTPUT_FILE_PATH` with the secrets in it

## How to use

### Build

```bash
docker-compose build
```

### Kubernetes

See [manifest.yaml](./k8s_templates/manifest.yaml) for an example on how to set it up
on kubernetes.

NOTE: you need to be using [Azure Active Directory pod-managed identities](https://github.com/Azure/aad-pod-identity). This is what permits the init container in the pod to connect to the Azure Key Vault.

You can pass azure credentials via these env vars if you want:

```text
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID
```

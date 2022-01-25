# docker-rest-puller-restarter

KISS REST API to pull and restart Docker containers

## Links

- GitHub: https://github.com/JenswBE/docker-rest-puller-restarter
- DockerHub: https://hub.docker.com/r/jenswbe/docker-rest-puller-restarter

## Configuration

Create a `config.yml` file:

```yaml
Clients: # Mandatory
  - Name: Test
    APIKey: TEST_API_KEY
    ContainerNames: nginx,traefik # Use * to allow all container names
Server: # Optional
  Debug: true # Defaults to false
  Port: 8080 # Defaults to 8080
  TrustedProxies: "172.16.0.0/16" # Defaults to 172.16.0.0/16, which is the default Docker IP range
```

See https://github.com/JenswBE/docker-rest-puller-restarter/blob/main/docker-compose.yml for a Docker Compose example.

## Usage

```bash
API_KEY=""
CONTAINER_NAME=""
curl --fail -vH "API-KEY: ${API_KEY:?}" -X POST "https://drpr-eve.jensw.be/${CONTAINER_NAME:?}/pull_restart/"
```

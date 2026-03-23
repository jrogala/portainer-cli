# portainer-cli

CLI for Docker container management via Portainer API.

## Setup

Set `PORTAINER_URL` env var or use config file, then authenticate:

```bash
export PORTAINER_URL=https://portainer.example.com
portainer login -u admin -p password
```

## Commands

| Command | Description |
|---|---|
| `login` | Authenticate with Portainer instance |
| `ps` | List running containers |
| `list` | List all containers (alias for ps) |
| `start` | Start a stopped container |
| `stop` | Stop a running container |
| `restart` | Restart a container |
| `logs` | Show container logs |
| `inspect` | Show detailed container info |
| `stacks` | List deployed stacks |

## Examples

```bash
# Login and list running containers
portainer login -u admin -p password
portainer ps

# Restart a container by name
portainer restart nginx-proxy

# View logs for a container
portainer logs --tail 100 my-app

# Inspect a container's full config
portainer inspect my-app
```

## JSON Output

All commands support `--json` for machine-readable output.

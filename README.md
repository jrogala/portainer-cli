# portainer-cli

CLI for Docker container management via Portainer API.

## Install

Download a binary from the [latest release](https://github.com/jrogala/portainer-cli/releases/latest), or install with Go:

```bash
go install github.com/jrogala/portainer-cli@latest
```

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
$ portainer ps
ID            NAME           IMAGE                                STATE    STATUS
9947ef62a516  vikunja        vikunja/vikunja:latest               running  Up 7 days
42a4c3b77f4b  jellyfin       jellyfin/jellyfin                    running  Up 5 days (healthy)
cc38151f9231  homeassistant  homeassistant/home-assistant:stable  running  Up 5 days
a7ae788e6885  portainer      portainer/portainer-ce:latest        running  Up 7 months

$ portainer stacks
ID  NAME              STATUS
1   homeassistant     active
2   factorio          active
22  jellyfin          active
25  vikunja           active

$ portainer restart jellyfin
Restarted jellyfin
```

## JSON Output

All commands support `--json` for machine-readable output.

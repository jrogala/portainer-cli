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
ID            NAME        IMAGE            STATE    STATUS
abcd12345678  nginx       nginx:latest     running  Up 2 days
efgh56789012  postgres    postgres:15      running  Up 5 hours

$ portainer restart nginx
Restarted nginx

$ portainer stacks
ID  NAME          STATUS
1   monitoring    active
2   media-stack   active

$ portainer inspect nginx
Name:     nginx
ID:       abcd12345678
Image:    nginx:latest
State:    running
Running:  true
Started:  2026-03-22T10:00:00Z
Restart:  always
Network:  bridge (172.17.0.2)
```

## JSON Output

All commands support `--json` for machine-readable output.

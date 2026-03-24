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
ID          NAME        IMAGE            STATE    STATUS
abcd1234    nginx       nginx:latest     running  Up 2 days
efgh5678    postgres    postgres:15      running  Up 5 hours
ijkl9012    backup      backup:1.0       exited   Exited (0) 3h ago

$ portainer restart nginx
Restarted nginx

$ portainer logs --tail 3 nginx
[nginx] 10:15:32 "GET / HTTP/1.1" 200
[nginx] 10:16:15 "GET /api HTTP/1.1" 200
[nginx] 10:17:01 "POST /api HTTP/1.1" 201

$ portainer inspect nginx
Name:     nginx
ID:       abcd1234
Image:    nginx:latest
State:    running
Started:  2026-03-22T10:00:00Z
Restart:  always
Network:  bridge (172.17.0.2)
```

## JSON Output

All commands support `--json` for machine-readable output.

// Package ops contains business logic for Portainer operations.
// Functions return Go structs and errors — zero I/O, zero formatting.
package ops

import (
	"fmt"

	"github.com/jrogala/portainer-cli/client"
)

// ContainerEntry represents a container in listing results.
type ContainerEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	State  string `json:"state"`
	Status string `json:"status"`
}

// ListContainers returns containers, optionally including stopped ones.
func ListContainers(c *client.Client, all bool) ([]ContainerEntry, error) {
	containers, err := c.ListContainers(all)
	if err != nil {
		return nil, err
	}
	var entries []ContainerEntry
	for _, ct := range containers {
		entries = append(entries, ContainerEntry{
			ID:     ct.ID[:12],
			Name:   ct.Name(),
			Image:  ct.Image,
			State:  ct.State,
			Status: ct.Status,
		})
	}
	return entries, nil
}

// StartContainer resolves and starts a container by name or ID.
func StartContainer(c *client.Client, name string) error {
	id, err := ResolveContainer(c, name)
	if err != nil {
		return err
	}
	return c.StartContainer(id)
}

// StopContainer resolves and stops a container by name or ID.
func StopContainer(c *client.Client, name string) error {
	id, err := ResolveContainer(c, name)
	if err != nil {
		return err
	}
	return c.StopContainer(id)
}

// RestartContainer resolves and restarts a container by name or ID.
func RestartContainer(c *client.Client, name string) error {
	id, err := ResolveContainer(c, name)
	if err != nil {
		return err
	}
	return c.RestartContainer(id)
}

// GetLogs returns container logs.
func GetLogs(c *client.Client, name string, tail int) (string, error) {
	id, err := ResolveContainer(c, name)
	if err != nil {
		return "", err
	}
	return c.GetLogs(id, tail)
}

// ContainerDetail represents detailed container info.
type ContainerDetail struct {
	Name          string            `json:"name"`
	ID            string            `json:"id"`
	Image         string            `json:"image"`
	State         string            `json:"state"`
	Running       bool              `json:"running"`
	StartedAt     string            `json:"started_at"`
	RestartPolicy string            `json:"restart_policy"`
	Networks      map[string]string `json:"networks"`
}

// InspectContainer returns detailed info about a container.
func InspectContainer(c *client.Client, name string) (*ContainerDetail, error) {
	id, err := ResolveContainer(c, name)
	if err != nil {
		return nil, err
	}
	info, err := c.InspectContainer(id)
	if err != nil {
		return nil, err
	}
	networks := make(map[string]string)
	for netName, net := range info.NetworkSettings.Networks {
		networks[netName] = net.IPAddress
	}
	return &ContainerDetail{
		Name:          info.Name,
		ID:            info.ID[:12],
		Image:         info.Config.Image,
		State:         info.State.Status,
		Running:       info.State.Running,
		StartedAt:     info.State.StartedAt,
		RestartPolicy: info.HostConfig.RestartPolicy.Name,
		Networks:      networks,
	}, nil
}

// ResolveContainer finds a container by name or ID prefix.
func ResolveContainer(c *client.Client, name string) (string, error) {
	containers, err := c.ListContainers(true)
	if err != nil {
		return "", err
	}
	for _, ct := range containers {
		if ct.Name() == name || ct.ID[:12] == name || ct.ID == name {
			return ct.ID, nil
		}
	}
	return "", fmt.Errorf("container %q not found", name)
}

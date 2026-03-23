package ops

import "github.com/jrogala/portainer-cli/client"

// StackEntry represents a stack in listing results.
type StackEntry struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ListStacks returns all stacks.
func ListStacks(c *client.Client) ([]StackEntry, error) {
	stacks, err := c.ListStacks()
	if err != nil {
		return nil, err
	}
	var entries []StackEntry
	for _, s := range stacks {
		status := "inactive"
		if s.Status == 1 {
			status = "active"
		}
		entries = append(entries, StackEntry{
			ID:     s.ID,
			Name:   s.Name,
			Status: status,
		})
	}
	return entries, nil
}

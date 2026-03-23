package tests

import (
	"fmt"
	"net/http"

	"github.com/jrogala/portainer-cli/client"
	"github.com/jrogala/portainer-cli/pkg/ops"
)

func (sc *scenarioCtx) aStackExistsWithStatusActive(name string) error {
	sc.stacks = append(sc.stacks, mockStack{
		ID:     len(sc.stacks) + 1,
		Name:   name,
		Status: 1,
	})
	sc.registerStackRoutes()
	return nil
}

func (sc *scenarioCtx) aStackExistsWithStatusInactive(name string) error {
	sc.stacks = append(sc.stacks, mockStack{
		ID:     len(sc.stacks) + 1,
		Name:   name,
		Status: 0,
	})
	sc.registerStackRoutes()
	return nil
}

func (sc *scenarioCtx) noStacksExist() error {
	sc.stacks = nil
	sc.registerStackRoutes()
	return nil
}

func (sc *scenarioCtx) registerStackRoutes() {
	var stacks []client.Stack
	for _, ms := range sc.stacks {
		stacks = append(stacks, client.Stack{
			ID:         ms.ID,
			Name:       ms.Name,
			Status:     ms.Status,
			EndpointID: 2,
		})
	}
	if stacks == nil {
		stacks = []client.Stack{}
	}
	sc.mock.On("GET", "/api/stacks", http.StatusOK, stacks)
}

func (sc *scenarioCtx) iListStacks() error {
	entries, err := ops.ListStacks(sc.client)
	sc.lastErr = err
	sc.stackList = entries
	return nil
}

func (sc *scenarioCtx) stackListShouldInclude(name, status string) error {
	entries, ok := sc.stackList.([]ops.StackEntry)
	if !ok {
		return fmt.Errorf("no stack list available")
	}
	for _, e := range entries {
		if e.Name == name && e.Status == status {
			return nil
		}
	}
	return fmt.Errorf("stack %q with status %q not found in list", name, status)
}

func (sc *scenarioCtx) stackListShouldBeEmpty() error {
	entries, ok := sc.stackList.([]ops.StackEntry)
	if ok && len(entries) > 0 {
		return fmt.Errorf("expected empty stack list, got %d entries", len(entries))
	}
	return nil
}

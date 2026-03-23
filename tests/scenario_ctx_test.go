package tests

import (
	"github.com/jrogala/portainer-cli/client"
)

// scenarioCtx holds per-scenario state.
type scenarioCtx struct {
	mock   *MockServer
	client *client.Client

	// per-scenario data
	lastErr        error
	authToken      string
	containerList  any
	containerDetail any
	logs           string
	stackList      any

	// mock data tracking
	containers map[string]mockContainer
	stacks     []mockStack
}

type mockContainer struct {
	ID    string
	Name  string
	Image string
	State string
	Logs  string
}

type mockStack struct {
	ID     int
	Name   string
	Status int
}

func newScenarioCtx(mock *MockServer) *scenarioCtx {
	return &scenarioCtx{
		mock:       mock,
		containers: make(map[string]mockContainer),
	}
}

// newClient creates a client pointing at the mock server.
func (sc *scenarioCtx) newClient() *client.Client {
	return client.NewWithBaseURL("test-token", 2, sc.mock.URL())
}

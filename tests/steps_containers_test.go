package tests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jrogala/portainer-cli/client"
	"github.com/jrogala/portainer-cli/pkg/ops"
)

func (sc *scenarioCtx) aContainerExistsWithState(name, state string) error {
	id := shortID()
	sc.containers[name] = mockContainer{
		ID:    id,
		Name:  name,
		Image: name + ":latest",
		State: state,
	}
	sc.registerContainerRoutes()
	return nil
}

func (sc *scenarioCtx) containerHasLogs(name, logs string) error {
	mc, ok := sc.containers[name]
	if !ok {
		return fmt.Errorf("container %q not set up", name)
	}
	mc.Logs = strings.ReplaceAll(logs, "\\n", "\n")
	sc.containers[name] = mc

	// Build Docker log stream frames (8-byte header per line)
	logLines := strings.Split(mc.Logs, "\n")
	var framed []byte
	for _, line := range logLines {
		lineBytes := []byte(line + "\n")
		header := make([]byte, 8)
		header[0] = 1 // stdout
		size := len(lineBytes)
		header[4] = byte(size >> 24)
		header[5] = byte(size >> 16)
		header[6] = byte(size >> 8)
		header[7] = byte(size)
		framed = append(framed, header...)
		framed = append(framed, lineBytes...)
	}

	sc.mock.On("GET", fmt.Sprintf("/api/endpoints/2/docker/containers/%s/logs", mc.ID),
		http.StatusOK, framed)
	return nil
}

func (sc *scenarioCtx) registerContainerRoutes() {
	// Build the container list for the mock
	running := make([]client.Container, 0)
	all := make([]client.Container, 0)

	for _, mc := range sc.containers {
		ct := client.Container{
			ID:    mc.ID,
			Names: []string{"/" + mc.Name},
			Image: mc.Image,
			State: mc.State,
		}
		all = append(all, ct)
		if mc.State == "running" {
			running = append(running, ct)
		}
	}

	// Register routes for listing (running only and all)
	sc.mock.On("GET", "/api/endpoints/2/docker/containers/json?all=false", http.StatusOK, running)
	sc.mock.On("GET", "/api/endpoints/2/docker/containers/json?all=true", http.StatusOK, all)

	// Register routes for each container's actions
	for _, mc := range sc.containers {
		sc.mock.On("POST", fmt.Sprintf("/api/endpoints/2/docker/containers/%s/start", mc.ID),
			http.StatusNoContent, nil)
		sc.mock.On("POST", fmt.Sprintf("/api/endpoints/2/docker/containers/%s/stop", mc.ID),
			http.StatusNoContent, nil)
		sc.mock.On("POST", fmt.Sprintf("/api/endpoints/2/docker/containers/%s/restart", mc.ID),
			http.StatusNoContent, nil)

		// Register inspect route
		inspect := client.ContainerInspect{
			ID:   mc.ID,
			Name: "/" + mc.Name,
		}
		inspect.State.Status = mc.State
		inspect.State.Running = mc.State == "running"
		inspect.State.StartedAt = "2026-03-20T10:00:00Z"
		inspect.Config.Image = mc.Image
		inspect.HostConfig.RestartPolicy.Name = "always"
		inspect.NetworkSettings.Networks = map[string]struct {
			IPAddress string `json:"IPAddress"`
		}{
			"bridge": {IPAddress: "172.17.0.2"},
		}
		sc.mock.On("GET", fmt.Sprintf("/api/endpoints/2/docker/containers/%s/json", mc.ID),
			http.StatusOK, inspect)
	}
}

func (sc *scenarioCtx) iListContainers() error {
	entries, err := ops.ListContainers(sc.client, false)
	sc.lastErr = err
	sc.containerList = entries
	return nil
}

func (sc *scenarioCtx) iListAllContainers() error {
	entries, err := ops.ListContainers(sc.client, true)
	sc.lastErr = err
	sc.containerList = entries
	return nil
}

func (sc *scenarioCtx) containerListShouldInclude(name string) error {
	entries, ok := sc.containerList.([]ops.ContainerEntry)
	if !ok {
		return fmt.Errorf("no container list available")
	}
	for _, e := range entries {
		if e.Name == name {
			return nil
		}
	}
	return fmt.Errorf("container %q not found in list", name)
}

func (sc *scenarioCtx) iStartContainer(name string) error {
	sc.lastErr = ops.StartContainer(sc.client, name)
	return nil
}

func (sc *scenarioCtx) iStopContainer(name string) error {
	sc.lastErr = ops.StopContainer(sc.client, name)
	return nil
}

func (sc *scenarioCtx) iRestartContainer(name string) error {
	sc.lastErr = ops.RestartContainer(sc.client, name)
	return nil
}

func (sc *scenarioCtx) iGetLogsForContainer(name string) error {
	logs, err := ops.GetLogs(sc.client, name, 100)
	sc.lastErr = err
	sc.logs = logs
	return nil
}

func (sc *scenarioCtx) logsShouldContain(text string) error {
	if sc.lastErr != nil {
		return fmt.Errorf("expected logs but got error: %v", sc.lastErr)
	}
	if !strings.Contains(sc.logs, text) {
		return fmt.Errorf("logs do not contain %q, got: %q", text, sc.logs)
	}
	return nil
}

func (sc *scenarioCtx) iInspectContainer(name string) error {
	detail, err := ops.InspectContainer(sc.client, name)
	sc.lastErr = err
	sc.containerDetail = detail
	return nil
}

func (sc *scenarioCtx) inspectionShouldShowName(name string) error {
	detail, ok := sc.containerDetail.(*ops.ContainerDetail)
	if !ok || detail == nil {
		return fmt.Errorf("no inspection result available")
	}
	// Container names from Docker have leading slash
	cleanName := strings.TrimPrefix(detail.Name, "/")
	if cleanName != name {
		return fmt.Errorf("expected name %q, got %q", name, cleanName)
	}
	return nil
}

func (sc *scenarioCtx) inspectionShouldShowState(state string) error {
	detail, ok := sc.containerDetail.(*ops.ContainerDetail)
	if !ok || detail == nil {
		return fmt.Errorf("no inspection result available")
	}
	if detail.State != state {
		return fmt.Errorf("expected state %q, got %q", state, detail.State)
	}
	return nil
}

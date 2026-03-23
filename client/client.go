package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Client communicates with the Portainer API.
type Client struct {
	baseURL    string
	token      string
	endpointID int
	httpClient *http.Client
}

// New creates a new Portainer API client with TLS skip verify.
func New(baseURL, token string, endpointID int) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		endpointID: endpointID,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// NewWithBaseURL creates a client for testing — no TLS, custom base URL.
func NewWithBaseURL(token string, endpointID int, baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		endpointID: endpointID,
		httpClient: &http.Client{},
	}
}

func (c *Client) do(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("portainer error %d: %s", resp.StatusCode, string(data))
	}

	return resp, nil
}

func (c *Client) dockerPath(path string) string {
	return "/api/endpoints/" + strconv.Itoa(c.endpointID) + "/docker" + path
}

// EndpointID returns the configured endpoint ID.
func (c *Client) EndpointID() int {
	return c.endpointID
}

// Authenticate gets a JWT token from Portainer.
func Authenticate(baseURL, username, password string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	resp, err := client.Post(baseURL+"/api/auth", "application/json", strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed (%d): %s", resp.StatusCode, string(data))
	}

	var result struct {
		JWT string `json:"jwt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.JWT, nil
}

// ListContainers returns all containers.
func (c *Client) ListContainers(all bool) ([]Container, error) {
	path := c.dockerPath("/containers/json?all=" + strconv.FormatBool(all))
	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []Container
	return containers, json.NewDecoder(resp.Body).Decode(&containers)
}

// StartContainer starts a container by ID or name.
func (c *Client) StartContainer(id string) error {
	resp, err := c.do("POST", c.dockerPath("/containers/"+id+"/start"), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// StopContainer stops a container by ID or name.
func (c *Client) StopContainer(id string) error {
	resp, err := c.do("POST", c.dockerPath("/containers/"+id+"/stop"), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// RestartContainer restarts a container by ID or name.
func (c *Client) RestartContainer(id string) error {
	resp, err := c.do("POST", c.dockerPath("/containers/"+id+"/restart"), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// GetLogs returns container logs.
func (c *Client) GetLogs(id string, tail int) (string, error) {
	path := c.dockerPath("/containers/" + id + "/logs?stdout=true&stderr=true&tail=" + strconv.Itoa(tail))
	resp, err := c.do("GET", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Strip Docker log header bytes (8 bytes per line)
	return stripLogHeaders(data), nil
}

// InspectContainer returns detailed container info.
func (c *Client) InspectContainer(id string) (*ContainerInspect, error) {
	resp, err := c.do("GET", c.dockerPath("/containers/"+id+"/json"), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info ContainerInspect
	return &info, json.NewDecoder(resp.Body).Decode(&info)
}

// ListStacks returns all stacks.
func (c *Client) ListStacks() ([]Stack, error) {
	resp, err := c.do("GET", "/api/stacks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stacks []Stack
	return stacks, json.NewDecoder(resp.Body).Decode(&stacks)
}

// stripLogHeaders removes the 8-byte Docker log stream header from each line.
func stripLogHeaders(data []byte) string {
	var lines []string
	for len(data) > 0 {
		if len(data) < 8 {
			break
		}
		// Read size from header bytes 4-7
		size := int(data[4])<<24 | int(data[5])<<16 | int(data[6])<<8 | int(data[7])
		data = data[8:]
		if size > len(data) {
			size = len(data)
		}
		line := strings.TrimRight(string(data[:size]), "\n")
		if line != "" {
			lines = append(lines, line)
		}
		data = data[size:]
	}
	return strings.Join(lines, "\n")
}

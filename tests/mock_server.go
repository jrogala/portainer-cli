package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// MockRoute represents a registered mock response.
type MockRoute struct {
	Method string
	Path   string
	Status int
	Body   any
}

// MockServer is a test HTTP server that returns pre-registered responses.
type MockServer struct {
	Server *httptest.Server
	mu     sync.Mutex
	routes []MockRoute
}

// NewMockServer creates and starts a mock HTTP server.
func NewMockServer() *MockServer {
	ms := &MockServer{}
	ms.Server = httptest.NewServer(http.HandlerFunc(ms.handler))
	return ms
}

// On registers a response for a method+path combination.
// If the same method+path was already registered, it is replaced.
func (ms *MockServer) On(method, path string, status int, body any) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	// Replace existing route if found
	for i, r := range ms.routes {
		if r.Method == method && r.Path == path {
			ms.routes[i] = MockRoute{Method: method, Path: path, Status: status, Body: body}
			return
		}
	}
	ms.routes = append(ms.routes, MockRoute{
		Method: method,
		Path:   path,
		Status: status,
		Body:   body,
	})
}

// Reset clears all registered routes.
func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.routes = nil
}

// Close shuts down the mock server.
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// URL returns the mock server's URL.
func (ms *MockServer) URL() string {
	return ms.Server.URL
}

func (ms *MockServer) handler(w http.ResponseWriter, r *http.Request) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	path := r.URL.Path
	// Include query string for matching if present
	if r.URL.RawQuery != "" {
		path = path + "?" + r.URL.RawQuery
	}

	for _, route := range ms.routes {
		if route.Method == r.Method && matchPath(route.Path, path) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(route.Status)
			if route.Body != nil {
				switch v := route.Body.(type) {
				case string:
					w.Write([]byte(v))
				case []byte:
					w.Write(v)
				default:
					json.NewEncoder(w).Encode(v)
				}
			}
			return
		}
	}

	// No match found
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "no mock registered for " + r.Method + " " + path,
	})
}

// matchPath checks if a registered path matches the request path.
// Supports prefix matching with trailing wildcard.
func matchPath(pattern, path string) bool {
	if pattern == path {
		return true
	}
	// Also match ignoring query string differences
	if strings.Contains(pattern, "?") {
		return pattern == path
	}
	// Match pattern without query against path without query
	pathNoQuery := strings.SplitN(path, "?", 2)[0]
	return pattern == pathNoQuery
}

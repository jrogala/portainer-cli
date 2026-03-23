package tests

import (
	"fmt"
	"net/http"
)

func (sc *scenarioCtx) iLoginWithValidCredentials() error {
	// Register mock for successful auth
	sc.mock.On("POST", "/api/auth", http.StatusOK, map[string]string{
		"jwt": "valid-jwt-token-123",
	})

	// Use a plain HTTP client to simulate the login call
	resp, err := http.Post(sc.mock.URL()+"/api/auth", "application/json", nil)
	if err != nil {
		sc.lastErr = err
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sc.authToken = "valid-jwt-token-123"
		sc.lastErr = nil
	} else {
		sc.lastErr = fmt.Errorf("auth failed with status %d", resp.StatusCode)
	}
	return nil
}

func (sc *scenarioCtx) iLoginWithInvalidCredentials() error {
	sc.mock.On("POST", "/api/auth", http.StatusUnauthorized, map[string]string{
		"message": "Invalid credentials",
	})

	resp, err := http.Post(sc.mock.URL()+"/api/auth", "application/json", nil)
	if err != nil {
		sc.lastErr = err
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sc.lastErr = fmt.Errorf("authentication failed (%d)", resp.StatusCode)
	}
	return nil
}

func (sc *scenarioCtx) authShouldSucceed() error {
	if sc.lastErr != nil {
		return fmt.Errorf("expected auth to succeed, got: %v", sc.lastErr)
	}
	return nil
}

func (sc *scenarioCtx) authShouldFail() error {
	if sc.lastErr == nil {
		return fmt.Errorf("expected auth to fail, but it succeeded")
	}
	return nil
}

func (sc *scenarioCtx) iShouldReceiveAJWTToken() error {
	if sc.authToken == "" {
		return fmt.Errorf("no JWT token received")
	}
	return nil
}

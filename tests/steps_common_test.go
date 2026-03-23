package tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func initializeScenario(ctx *godog.ScenarioContext) {
	sc := newScenarioCtx(globalMock)

	ctx.Before(func(ctx context.Context, sc2 *godog.Scenario) (context.Context, error) {
		sc.mock.Reset()
		sc.lastErr = nil
		sc.authToken = ""
		sc.containerList = nil
		sc.containerDetail = nil
		sc.logs = ""
		sc.stackList = nil
		sc.containers = make(map[string]mockContainer)
		sc.stacks = nil
		sc.client = sc.newClient()
		// Register empty default responses
		sc.registerContainerRoutes()
		return ctx, nil
	})

	// --- Background ---
	ctx.Step(`^a running Portainer instance$`, sc.aRunningPortainerInstance)
	ctx.Step(`^an authenticated user$`, sc.anAuthenticatedUser)

	// --- Auth ---
	ctx.Step(`^I login with valid credentials$`, sc.iLoginWithValidCredentials)
	ctx.Step(`^I login with invalid credentials$`, sc.iLoginWithInvalidCredentials)
	ctx.Step(`^authentication should succeed$`, sc.authShouldSucceed)
	ctx.Step(`^authentication should fail$`, sc.authShouldFail)
	ctx.Step(`^I should receive a JWT token$`, sc.iShouldReceiveAJWTToken)

	// --- Containers ---
	ctx.Step(`^a container "([^"]*)" exists with state "([^"]*)"$`, sc.aContainerExistsWithState)
	ctx.Step(`^container "([^"]*)" has logs "([^"]*)"$`, sc.containerHasLogs)
	ctx.Step(`^I list containers$`, sc.iListContainers)
	ctx.Step(`^I list all containers$`, sc.iListAllContainers)
	ctx.Step(`^the container list should include "([^"]*)"$`, sc.containerListShouldInclude)
	ctx.Step(`^I start container "([^"]*)"$`, sc.iStartContainer)
	ctx.Step(`^I stop container "([^"]*)"$`, sc.iStopContainer)
	ctx.Step(`^I restart container "([^"]*)"$`, sc.iRestartContainer)
	ctx.Step(`^the operation should succeed$`, sc.operationShouldSucceed)
	ctx.Step(`^the operation should fail with "([^"]*)"$`, sc.operationShouldFailWith)
	ctx.Step(`^I get logs for container "([^"]*)"$`, sc.iGetLogsForContainer)
	ctx.Step(`^the logs should contain "([^"]*)"$`, sc.logsShouldContain)
	ctx.Step(`^I inspect container "([^"]*)"$`, sc.iInspectContainer)
	ctx.Step(`^the inspection should show name "([^"]*)"$`, sc.inspectionShouldShowName)
	ctx.Step(`^the inspection should show state "([^"]*)"$`, sc.inspectionShouldShowState)

	// --- Stacks ---
	ctx.Step(`^a stack "([^"]*)" exists with status active$`, sc.aStackExistsWithStatusActive)
	ctx.Step(`^a stack "([^"]*)" exists with status inactive$`, sc.aStackExistsWithStatusInactive)
	ctx.Step(`^no stacks exist$`, sc.noStacksExist)
	ctx.Step(`^I list stacks$`, sc.iListStacks)
	ctx.Step(`^the stack list should include "([^"]*)" with status "([^"]*)"$`, sc.stackListShouldInclude)
	ctx.Step(`^the stack list should be empty$`, sc.stackListShouldBeEmpty)
}

// shortID returns a random hex string for unique IDs.
func shortID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// --- Background steps ---

func (sc *scenarioCtx) aRunningPortainerInstance() error {
	if sc.mock == nil {
		return fmt.Errorf("mock server is not running")
	}
	return nil
}

func (sc *scenarioCtx) anAuthenticatedUser() error {
	if sc.client == nil {
		return fmt.Errorf("no authenticated client")
	}
	return nil
}

// --- Common assertion ---

func (sc *scenarioCtx) operationShouldSucceed() error {
	if sc.lastErr != nil {
		return fmt.Errorf("expected operation to succeed, got: %v", sc.lastErr)
	}
	return nil
}

func (sc *scenarioCtx) operationShouldFailWith(expected string) error {
	if sc.lastErr == nil {
		return fmt.Errorf("expected an error containing %q but got none", expected)
	}
	if !strings.Contains(sc.lastErr.Error(), expected) {
		return fmt.Errorf("expected error containing %q, got: %s", expected, sc.lastErr.Error())
	}
	return nil
}

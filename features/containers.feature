Feature: Container management
  As a user I can list, start, stop, restart, inspect, and view logs of containers.

  Background:
    Given a running Portainer instance
    And an authenticated user

  Scenario: List running containers
    Given a container "web-app" exists with state "running"
    And a container "db-server" exists with state "running"
    When I list containers
    Then the container list should include "web-app"
    And the container list should include "db-server"

  Scenario: List all containers including stopped
    Given a container "web-app" exists with state "running"
    And a container "old-task" exists with state "exited"
    When I list all containers
    Then the container list should include "web-app"
    And the container list should include "old-task"

  Scenario: Start a container
    Given a container "my-service" exists with state "exited"
    When I start container "my-service"
    Then the operation should succeed

  Scenario: Stop a container
    Given a container "my-service" exists with state "running"
    When I stop container "my-service"
    Then the operation should succeed

  Scenario: Restart a container
    Given a container "my-service" exists with state "running"
    When I restart container "my-service"
    Then the operation should succeed

  Scenario: Get container logs
    Given a container "my-service" exists with state "running"
    And container "my-service" has logs "Starting server...\nListening on port 8080"
    When I get logs for container "my-service"
    Then the logs should contain "Starting server"

  Scenario: Inspect a container
    Given a container "my-service" exists with state "running"
    When I inspect container "my-service"
    Then the inspection should show name "my-service"
    And the inspection should show state "running"

  Scenario: Resolve container not found
    When I start container "nonexistent"
    Then the operation should fail with "not found"

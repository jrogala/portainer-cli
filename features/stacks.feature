Feature: Stack management
  As a user I can list stacks and see their status.

  Background:
    Given a running Portainer instance
    And an authenticated user

  Scenario: List stacks
    Given a stack "production" exists with status active
    And a stack "staging" exists with status inactive
    When I list stacks
    Then the stack list should include "production" with status "active"
    And the stack list should include "staging" with status "inactive"

  Scenario: Empty stack list
    Given no stacks exist
    When I list stacks
    Then the stack list should be empty

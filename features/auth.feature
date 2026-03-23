Feature: Authentication
  As a user I can authenticate with Portainer and manage my session.

  Background:
    Given a running Portainer instance

  Scenario: Login with valid credentials
    When I login with valid credentials
    Then authentication should succeed
    And I should receive a JWT token

  Scenario: Login with invalid credentials
    When I login with invalid credentials
    Then authentication should fail

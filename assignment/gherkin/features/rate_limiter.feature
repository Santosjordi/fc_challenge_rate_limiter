Feature: Rate limiting of requests based on IP address and API token
  The server should limit the number of UUID generations per second
  based on either the client's IP address or an access token.
  The limit by token should override the IP-based limit.

  Background:
    Given the system is configured with:
      | Type         | Limit | Lockout Duration |
      | IP Default   | 5     | 5 minutes        |
      | Token abc123 | 10    | 2 minutes        |
    And the server is running and reachable on port 8080

  Scenario: Exceeding the IP limit results in HTTP 429
    Given a client with IP "192.168.1.1"
    When the client sends 6 requests within 1 second without an API token
    Then the response status code should be 429 for the 6th request
    And the response body should contain "you have reached the maximum number of requests or actions allowed"
    And the client should be blocked for 5 minutes

  Scenario: Token limit overrides IP-based limit
    Given a client with IP "192.168.1.1" and token "abc123"
    When the client sends 11 requests within 1 second with the token in the API_KEY header
    Then the response status code should be 429 for the 11th request
    And the client should be blocked by token rules for 2 minutes

  Scenario: A second token has a different, lower limit
    Given the system is configured with:
      | Type         | Limit | Lockout Duration |
      | Token xyz789 | 3     | 1 minute         |
    And a client with IP "10.0.0.2" and token "xyz789"
    When the client sends 4 requests within 1 second
    Then the response status code should be 429 for the 4th request
    And the client should be blocked for 1 minute

  Scenario: A blocked client is immediately rejected during lockout
    Given a client with IP "192.168.1.1" has already been blocked
    When the client sends another request during the lockout period
    Then the response status code should be 429
    And the response body should mention the client is blocked

  Scenario: A client without a token is limited by IP
    Given a client without an API token and IP "172.16.0.5"
    When the client sends 5 requests in 1 second
    Then all requests succeed
    When the client sends a 6th request
    Then the response status code should be 429

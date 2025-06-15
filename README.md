# ⚙️ Fullcycle Technical challenge -> Rate Limiter: How It Works & How to Configure It

## UUID Generator

#### 1. **Setup**
The server of this project has a simple endpoint `\generate` which returns an uuid.

- Ensure you have Docker and Docker Compose installed.
- From the project root, start the required services (including Redis) with:
  ```bash
  docker-compose up -d
  ```
- Build and run the server (if not using Docker Compose for the app):
  ```bash
  go run ./cmd/server/uuid-generator-server.go
  ```
  The server will listen on port `8080` by default.

#### 2. **Configuration**

- All rate limiter settings are controlled via the `.env` file in the project root.
- You can adjust limits and lockout durations for both IP and token-based requests by editing this file before starting the server.

#### 3. **Testing the Endpoints**

- The main endpoint is:
  ```
  GET /generate
  ```
  This returns a UUID if the request is allowed.

- **Test IP-based limiting:**
  - Use a tool like `ab` (ApacheBench) or `curl` in a loop to send more requests per second than allowed by `IP_LIMIT_PER_SECOND`.
  - Example:
    ```bash
    ab -n 20 -c 5 http://localhost:8080/generate
    ```
  - After exceeding the limit, you should receive HTTP 429 responses with the message:
    ```
    you have reached the maximum number of requests or actions allowed within a certain time frame
    ```

- **Test Token-based limiting:**
  - Add the `API_KEY` header to your requests.
    ```bash
    ab -n 20 -c 5 -H "API_KEY: mytoken123" http://localhost:8080/generate
    ```
  - The token limit (`TOKEN_LIMIT_PER_SECOND`) will override the IP limit for requests with a valid token.

- **Test Lockout:**
  - After exceeding the limit, further requests should be blocked for the configured lockout duration (`IP_LOCKOUT_SECONDS` or `TOKEN_LOCKOUT_SECONDS`).
  - After the lockout period, requests should be allowed again.


### 🚦 How the Rate Limiter Works

This project implements a **fixed-window rate limiter** using Redis.

* Each incoming request is identified by a **key**, which can be:

  * A **token** (`API_KEY:<value>`) if an `API_KEY` is present in the request header.
  * The **client IP address** otherwise.
* For each key, Redis keeps a counter of requests per **1-second window**.
* If the number of requests in that second exceeds the configured threshold:

  * The key is **locked out** for a defined duration.
  * Any request during lockout receives an HTTP **429 Too Many Requests** error.
* Rate-limiting headers are included in responses:

  * `X-RateLimit-Limit`
  * `X-RateLimit-Remaining`
  * `X-RateLimit-Reset`

---

### 🛠️ Configuration

All configuration is loaded from the `.env` file in the project root.

| Variable                 | Description                                         | Example     |
| ------------------------ | --------------------------------------------------- | ----------- |
| `SERVER_PORT`            | HTTP server port                                    | `:8080`     |
| `REDIS_HOST`             | Redis server address                                | `localhost` |
| `REDIS_PORT`             | Redis port                                          | `6379`      |
| `REDIS_PASSWORD`         | Redis password (optional)                           | \`\`        |
| `REDIS_DB`               | Redis DB index                                      | `0`         |
| `IP_LIMIT_PER_SECOND`    | Max requests per second for IP-based keys           | `5`         |
| `IP_LOCKOUT_SECONDS`     | Lockout duration (in seconds) for IP-based abuse    | `10`        |
| `TOKEN_LIMIT_PER_SECOND` | Max requests per second for token-based keys        | `10`        |
| `TOKEN_LOCKOUT_SECONDS`  | Lockout duration (in seconds) for token-based abuse | `5`         |

### 📂 Example `.env` File

```env
SERVER_PORT=:8080

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

IP_LIMIT_PER_SECOND=5
IP_LOCKOUT_SECONDS=10

TOKEN_LIMIT_PER_SECOND=10
TOKEN_LOCKOUT_SECONDS=5
```

---

### 🧪 Testing

You can test the rate limiter using tools like `ab` (ApacheBench):

```bash
ab -n 100 -c 10 -H "token:token:123" http://localhost:8080/generate
```

This sends 100 requests with concurrency 10, using a custom token header.

---

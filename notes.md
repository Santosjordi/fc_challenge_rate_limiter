## Configuration

This project uses a `.env` file for configuration, which should be placed at the root of the project directory. You can set up various environment variables in the `.env` file to customize the behavior of the rate limiter and other system parameters.

#### Required Configuration in `.env`

Below is an example configuration for the `.env` file:

```
# Rate limits
RATE_LIMIT_IP_DEFAULT=5                          # Rate limit per IP address (requests per second)
RATE_LIMIT_TOKENS=10                             # Rate limit per token (requests per second)

# Lockout durations (in seconds)
RATE_LIMIT_IP_LOCKOUT_DURATION=300               # Lockout duration for IP (in seconds)
RATE_LIMIT_TOKEN_LOCKOUT_DURATION=120            # Lockout duration for token (in seconds)

# Persistence backend configuration (default: Redis)
BACKEND=redis                                    # Choose 'redis' or 'memcached' as persistence backend
REDIS_HOST=localhost                             # Host for Redis (used when BACKEND=redis)
REDIS_PORT=6379                                  # Port for Redis (used when BACKEND=redis)

# JWT secret (for token validation)
JWT_SECRET=mysecretkey                           # JWT secret key used for token validation (if using JWT)

# Any other configuration parameters
```

#### How to Set Up the `.env` File

1. Copy the example configuration into a `.env` file at the root of the project.
2. Customize the values based on your environment or needs (e.g., changing the rate limits, JWT secret, etc.).
3. Make sure you have the required services (e.g., Redis) running if your configuration is using those backends.

#### Example Configuration

For example, if you're using Redis as the backend, the `.env` would look like this:

```
RATE_LIMIT_IP_DEFAULT=5
RATE_LIMIT_TOKENS=10
RATE_LIMIT_IP_LOCKOUT_DURATION=300
RATE_LIMIT_TOKEN_LOCKOUT_DURATION=120
BACKEND=redis
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=mysecretkey
```

If you're using a different backend (e.g., Memcached), just change the `BACKEND` variable accordingly, and add the configuration for that service (e.g., Memcached host and port).

---
### Running

Run the application from the project root:
```bash
go run cmd/server/uuid-generator-server.go
```


Use Apache Benchmark to load tests the application.

```bash
ab -n 100 -c 10 -H "API_KEY: mytoken" http://localhost:8080/generate

```

---
## Flow Diagram

```
Request → Middleware Chain
           ↓
    1. Extract Identifier (IP/Token)
           ↓
    2. Check Lockout Status
           ↓
    3. Get Rate Limit Config
           ↓
    4. Apply Rate Limiting
           ↓
    5. Register Request
           ↓
    Next Handler/Response

```

Architecture drawing:
```
+------------------------+            +------------------------+
|     config.Config      |            |       Viper/.env       |
|  [SRP]                 |<-----------+  Loads configuration   |
+------------------------+            +------------------------+

         |
         v
+------------------------+            +-------------------------+
|     main.go            |            |    Redis DB             |
|  [SRP]                 |            |                         |
|  - Loads config        |            +-------------------------+
|  - Instantiates        |
|    limiter & middleware|
+-----------+------------+
            |
            v
+------------------------+           Implements
|   RedisLimiter         |------------------------------+
| [SRP, OCP, LSP]        |                              |
|  - Implements          |                              |
|    RateLimiter         |                              |
|  - Talks to Redis      |                              v
+-----------+------------+                  +----------------------+
            |                                |  RateLimiter        |
            | Uses                           |  Interface          |
            +------------------------------->|  [ISP, DIP]         |
                                             +----------------------+
                                                       ^
                                                       |
                             +-------------------------+
                             |
                             v
                    +--------------------------+
                    | RateLimitMiddleware      |
                    | [SRP, OCP, DIP]          |
                    | - Uses RateLimiter       |
                    | - HTTP 429 on limit hit  |
                    +-------------+------------+
                                  |
                                  v
                         +------------------+
                         |  HTTP Requests    |
                         |  (via chi router) |
                         +------------------+


```


---
## Project Structure

```
📦 fc_challenge_rate_limiter
├── 📁 cmd
│   └── 📁 server
│       └── 📄 uuid-generator-server.go
├── 📁 config
│   ├── 📄 README.md
│   ├── 📄 config.go
│   └── 📄 config_test.go
├── 📁 docs
│   ├── 📄 ASSIGNMENT.md
│   ├── 📄 architecture.excalidraw
│   └── 📁 gherkin
│       └── 📁 features
│           ├── 📄 rate_limiter.feature
│           └── 📄 rate_limiter_test.go
├── 📁 internal
│   ├── 📁 infra
│   │   └── 📁 db
│   │       ├── 📄 interface.go
│   │       └── 📄 redis_rate_limiter.go
│   └── 📁 webserver
│       ├── 📁 handlers
│       │   ├── 📄 uuid_handler.go
│       │   └── 📄 uuid_handler_test.go
│       ├── 📁 middleware
│       │   └── 📄 rate_limiter.go
│       └── 📁 utils
│           └── 📄 http_utils.go
├── 📁 pkg
│   └── 📁 entity
│       └── 📄 ID.go
├── 📁 test
│   └── 📁 testdata
│       └── 📄 .env.test
├── 📄 .env
├── 📄 .gitignore
├── 📄 README.md
├── 📄 docker-compose.yml
├── 📄 go.mod
└── 📄 go.sum
```


---

## Dependency Tree

```
Primary Dependencies:
├── cmd/server/uuid-generator-server.go
│   ├── config/config.go
│   ├── internal/webserver/handlers/uuid_handler.go
│   └── chi router (external)
│
├── internal/webserver/middleware/rate_limiter.go
│   ├── internal/infra/db/interface.go
│   ├── internal/webserver/utils/http_utils.go
│   └── net/http (stdlib)
│
├── internal/infra/db/redis_rate_limiter.go
│   ├── internal/infra/db/interface.go
│   └── go-redis/v9 (external)
│
├── config/config.go
│   ├── go-chi/jwtauth (external)
│   └── spf13/viper (external)
│
└── pkg/entity/ID.go
    └── google/uuid (external)

External Dependencies:
├── github.com/go-chi/chi v1.5.5
├── github.com/go-chi/jwtauth v1.2.0
├── github.com/google/uuid v1.6.0
├── github.com/redis/go-redis/v9 v9.7.3
└── github.com/spf13/viper v1.20.1

```


---



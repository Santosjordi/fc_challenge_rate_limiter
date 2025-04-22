# Fullcycle challenge: Rate Limiter

## Desafio técnico pós graduação em Golang

Anotações para documentar:
- usar viper para acessar as variáveis de ambiente no `.env`

---

### Configuration

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

### Additional Notes

- The `.env` file is loaded automatically at startup using the [Viper](https://github.com/spf13/viper) package.
- You can override the settings in `.env` by setting the environment variables directly in your shell.

---

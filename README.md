# Microservice nakup

Provides timeslot reservation and purchase management.

## Env vars

Check out .env.example for example values

| ENV                         | Description                          |
| --------------------------- | ------------------------------------ |
| LOG_LEVEL                   | Log level (DEBUG, INFO, WARN, ERROR) |
| TZ                          | Timezone                             |
| POSTGRES_IP                 | Postgres DB IP                       |
| POSTGRES_PORT               | Postgres DB port                     |
| POSTGRES_USERNAME           | Postgres DB username                 |
| POSTGRES_PASSWORD           | Postgres DB password                 |
| POSTGRES_DATABASE_NAME      | Postgres DB database                 |
| POSTGRES_TEST_DATABASE_NAME | Postgres DB database for tests       |
| AUTH_HOST                   | Address of auth microservice         |
| SPORED_HOST                 | Address of spored microservice       |

## Running

Run the application via

```shell
godotenv go run main.go
```

Regenerate swagger docs via

```shell
make docs
```

Regenerate swagger clients via

```shell
make swagger-clients
```

Run all application tests via

```shell
make test
```

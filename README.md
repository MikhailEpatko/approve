![LoC Badge](https://github.com/MikhailEpatko/approve/blob/loc-badge/loc-badge.svg) ![HoC Badge](https://github.com/MikhailEpatko/approve/blob/hoc-badge/hoc-badge.svg)
# approve
### Approving Management

version: 0.0.0

1. Start PostgreSQL server

```shell
podman-compose -f ./container/podman-compose.yaml up
```

2. Create database

```shell
create database approve;
```

3. Add environment variables:

```bash
export DB_DRIVER_NAME=postgres
export DB_HOST=<host>
export DB_PORT=<port>
export DB_NAME=<db_name>
export DB_SSL_MODE=disable
export DB_USER=<user>
export DB_PASSWORD=<password>

export GOOSE_DRIVER=$DB_DRIVER_NAME
export GOOSE_DBSTRING=$DB_DRIVER_NAME://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME
export GOOSE_MIGRATION_DIR=./db/migrations

```
4. Install goose

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```
5. Run migrations:

```shell
goose up
```

![LoC Badge](https://github.com/MikhailEpatko/approve/blob/loc-badge/loc-badge.svg) ![HoC Badge](https://github.com/MikhailEpatko/approve/blob/hoc-badge/hoc-badge.svg)
# approve
### Approving Management


1. Start PostgreSQL server

```shell
podman-compose -f ./container/podman-compose.yaml up
```

2. Create database

```shell
create database approve;
```

3. Add .env file in the project root dir:

```dotenv
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://<user>:<password>@localhost:5432/approve
GOOSE_MIGRATION_DIR=./db/migrations
```
4. Install goose

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```
5. Run migrations:

```shell
goose up
```

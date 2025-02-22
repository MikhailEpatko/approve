# gidm
### Identity Management with Golang

1. Start PostgreSQL server

```shell
podman-compose -f ./container/podman-compose.yaml up
```

2. Create database

```shell
create database gidm;
```

3. Add .env file in the project root dir:

```dotenv
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://<user>:<password>@localhost:5432/gidm
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

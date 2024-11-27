# user_service

Created by:
Suprapto Sugiyanto (Prapto)

suprapto.sugiyanto@gmail.com

## Description

Here is a simple user application.

## Running

Here is the command to run the application:

```
go run cmd/api/main.go
```

## Mockgen

Here is the example to create mock:

```
mockgen -source=<source_file> -destination=<destination_file>
```

## migration

Ensure that you're alredy install golang-migrate.
Here is the command to create a migration file:

```
make migration name=<filename>
```

Here is the command to migrate:

```
make migrate url="postgres://<db-username>:<db-password>@localhost:5432/<db-name>"
```

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

mockgen -source=service/auth/auth.go -destination=test/mock/service/auth/auth.go
mockgen -source=service/user/user.go -destination=test/mock/service/user/user.go
mockgen -source=service/user/register/register.go -destination=test/mock/service/user/register/register.go
mockgen -source=service/user/login/login.go -destination=test/mock/service/user/login/login.go
```

## Migration

Ensure that you're alredy install golang-migrate.
Here is the command to create a migration file:

```
make migration name=<filename>
```

Here is the command to migrate:

```
make migrate url="postgres://<db-username>:<db-password>@localhost:5432/<db-name>"
```

## Endpoints

```
POST - /register
```

Request body

```
{
    "name":"Cristiano Ronaldo",
    "username":"ronaldo",
    "phone_number":"0856241",
    "password":"secret_agent"
}
```

Response

```
{
    "token": "asd123"
}
```

```
POST - /login
```

Request body

```
{
    "username":"ronaldo",
    "password":"secret_agent"
}
```

Response

```
{
    "token": "asd123"
}
```

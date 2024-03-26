# Snippetbox app

Simple app from Alex Edward's book "let's Go". Written in Go 1.22

## How to run
Just use docker compose, this will create MySQL and Go backend services. All Go files are watched by [air](https://github.com/cosmtrek/air): any changes will be detected and app will rebuild

```shell
docker compose up -d
```

App can be found on https://localhost:4000 (yes, it has self-signed certificate)

## Powered by
 - [scs](github.com/alexedwards/scs/) for session management
 - [httprouter](github.com/julienschmidt/httprouter) for routing
 - [alice](github.com/justinas/alice) for chaining middlewares
 - [nosurf](github.com/justinas/nosurf) CSRF
 - [crypto](golang.org/x/crypto) for password hashing

## Project structure
### cmd/web
Core logic of app
- handlers.go, routes.go – logic of handlers and URL patterns
- middleware.go – some middlewares, e.g. logging, authentication and permission checks
- templates.go – helpers functions for templates
- helpers.go - useful shortcuts, e.g. `serverError` and `render`

### internal
 - models representation, SQL
 - functions for testing
 - validation

### ui
HTML with Go templating


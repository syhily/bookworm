# Bookworm

This is a Golang written blog platform for [且听书吟](https://yufan.me).

This application is divided up into multiple packages, each with its own purpose.

- `components` - templ components.
- `db` - Database access code used to increment and get counts.
- `handlers` - HTTP handlers.
- `services` - Services used by the handlers.
- `sessions` - Middleware for implementing HTTP session IDs.
- `main.go` - Used to run the application locally.

## Development

[Templiér](https://github.com/romshark/templier) is used for local development.

# Outshort

**Simple link shortener web application**

![thumbnail](assets/home_page.png)

## About

### Stack

-   **Frontend framework:** Angular 19

-   **Styles:** PrimeNG 19 + Tailwind 3.4

-   **Server:** Go + Gin

-   **Database:** SQLite

## Quickstart

### Dependencies

-   `go`
-   `goose` ([go package](https://github.com/pressly/goose))
-   `docker` (recommended)

### Setup environment

Data inside `server/env` directory and `server/.env` file (TODO) are meant to be stored and modified outside the server container and mounted using volumes

```sh
$ cp server/.env.template server/.env

# TODO: additional configuration
```

### Initialize database

If you are running the project for the first time you need to initialize sqlite3 database and apply migrations using `goose`. This database will be stored locally and passed to the server container

```sh
$ cd server

# Run goose go initialize database
$ go run github.com/pressly/goose/v3/cmd/goose@latest up
```

This will generate a `server/env/database.db` file with the latest schema

### Build and run

```sh
# Setup permissions
$ sudo chmod +x ./scripts/local-deploy.sh

# Build client and server containers and run them
$ ./scripts/local-deploy.sh
```

This will create a copy of the development environment inside the `deploy` directory

## Development

**Start server**

```sh
$ cd server
$ go run .
```

**Start client**

```sh
$ cd client

# Install dependencies
$ npm i

$ npm start
```

After that open `http://localhost:4200` in your browser

## License

MIT

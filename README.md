# Brotchaschte

Brotchaschte is an order management and delivery planner application developed for the yearly occurring fundraiser ZoBa
of Cevi Oberi.

## Getting Started

### Building from Source

Install these prerequisites
- [Taskfile](https://taskfile.dev/)
- [Go >= 1.27](https://go.dev/)
- [pnpm](https://pnpm.io/)
- [sqlc](https://sqlc.dev/)
- and a C compiler like gcc as it uses CGO

Then, build a binary with
```sh
task build
```
which creates the `main` executable in the project root.

### Nix Flake
Brotchaschte is also packaged as a Nix Flake. So you can run and build the project using nix.
For example:
```sh
nix build github:theswampire/brotchaschte
```

### Running the Application
Running the server is as simple as
```sh
./main
```
The server will create an SQLite database `data.db` in the working directory and
start listening on port `8080` for connections.

The initial login credentials are
- **E-Mail**: admin
- **Password**: password123

Change them under "Einstellungen".


### Configuration
These are the configuration options for the server. Set these options as **environment variables**.

**PORT**. Set the port on which the server will listen. **Default:** `8080`.

**DATABASE_URL**. The database connection string to connect to an upstream PostgreSQL
database instance to import orders from. 

**ENABLE_RANDOM_ORDERS**. If this is set to `1`, then `/api/v1/orders/generate` is 
enabled which generates random orders *for debugging purposes*. **Default:** disabled.


## Development
### Dev-Shell
The Nix Flake provides a developer shell which has all the necessary dependencies already installed.
Use it with
```sh
nix develop
```

### Live / Dev Server
To use the live reload, install [air](https://github.com/air-verse/air).

Then run 
```sh
task
```
This builds and starts a version of the server where instead of serving static assets,
all non-api requests are proxied to a vite dev server.


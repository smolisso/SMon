# SMon

An open-source monitoring system by Stefano Molisso, built on experience with
Nagios since 2007. A modular Go monolith, PostgreSQL, and an embedded React /
TypeScript UI. Licensed under MIT.

## Current stage: foundation

Implemented:

- A Go HTTP server with graceful shutdown and a loopback default address.
- Embedded React UI displaying actual server and database readiness.
- `GET /api/v1/health`: process liveness.
- `GET /api/v1/ready`: database connectivity and migration readiness. Returns
  HTTP 503 if the database is missing, inaccessible, or not migrated.
- An explicit, transactional `smon migrate` command with a PostgreSQL advisory lock.
- Initial tables for hosts, commands, services, current service states and state events.

Not yet implemented: host/service management, plugin execution, macro expansion,
scheduling, notifications, authentication, graphs and configuration import.
Nagios compatibility is a goal, not a claim of this initial scaffold. The UI
shows no invented monitoring data and no plugins are executed.

## Build

Requirements: Go 1.25 or later, Node.js 22.12 or later with npm. PostgreSQL 16 or
later is the initial database target. Plugin interpreters are not needed yet.

```sh
npm --prefix web ci
npm --prefix web run build
go build -o bin/smon ./cmd/smon
```

Build the UI before any Go build or test: `web/dist` is embedded at compile time.
Node.js is needed only for development/builds, not on the deployment server.
Repository: https://github.com/smolisso/SMon.

## Run without a database

```sh
./bin/smon
```

Open http://127.0.0.1:8080. Liveness works, readiness returns 503 and the UI says
the database needs configuration. This mode does not persist data.

## Connect a dedicated PostgreSQL database

First provision a dedicated `smon` database and role using your usual PostgreSQL
administration process. Set `SMON_DATABASE_URL` in the environment. The example
below is for a local database; use the appropriate TLS settings for remote access.

```sh
export SMON_DATABASE_URL='postgres://smon:YOUR_PASSWORD@127.0.0.1:5432/smon?sslmode=disable'
./bin/smon migrate
./bin/smon
```

Migrations run only when requested and only against that database. Repeating
`migrate` skips applied migrations. `.env.example` documents the settings; the
application does not load `.env` files automatically.

The current server has no authentication. Keep the default loopback binding for
this development stage. `SMON_LISTEN` overrides `127.0.0.1:8080` explicitly.

## Verify

After building the frontend:

```sh
go test ./...
go vet ./...
```

To exercise migration idempotency and readiness on a **dedicated, empty test
database**, set `SMON_TEST_DATABASE_URL` before running tests. The integration
test creates schema objects there and leaves them for inspection; never point
it at a production or shared database.

## Structure and next step

See [architecture](docs/architecture.md) and the
[Nagios compatibility baseline](docs/nagios-compatibility.md).

The next milestone is a tested plugin runner using harmless local fixtures:
exit status, multiline output, performance data, timeouts and process cleanup.
Then connect scheduling, persistent results, and host/service management.

## Contributing

Keep changes small and document the supported behavior. Add meaningful tests for
monitoring semantics and execution boundaries. Never commit credentials, private
keys, production configurations or captured output containing secrets.

UI references: Stefano Molisso's MIT-licensed Nagios Modern UI. No files or
production settings have been copied from the inspected repositories.

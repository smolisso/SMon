# Architecture

One `smon` process serves the REST API and embedded UI and, in later milestones,
will schedule and execute checks. PostgreSQL holds configuration and state. No
message broker or separate frontend service is required.

| Module | Responsibility | Foundation status |
| --- | --- | --- |
| `cmd/smon` | Startup, configuration, shutdown, migration command | Implemented |
| `internal/api` | REST transport and embedded UI | Health/readiness only |
| `internal/postgres` | Persistence and transactional migration | Initial schema |
| `internal/monitoring` | Domain state, SOFT/HARD transitions and attempts | Boundary only |
| `internal/scheduler` | Due checks, retry timing, bounded concurrency | Boundary only |
| `internal/runner` | Process lifecycle, output limits and deadlines | Boundary only |
| `internal/nagios` | Output protocol, macros, configuration adapter | Boundary only |
| `web` | React/TypeScript, static asset embedding | Status overview |

The scheduler will ask the runner to execute a resolved command. The monitoring
module will interpret its result and persist current state plus state-change
events in one transaction. API writes must use those same domain operations.
Runner process failures must not be mistaken for healthy checks. Host states
need their own semantics; service exit-code mapping is not a host-state model.

The first schema models services only. A missing state row means pending, never
OK. Command arguments are arrays to preserve argument boundaries. Raw Nagios
command-line parsing and opt-in shell commands require an explicit adapter.

Before activating execution: define macro/quoting behavior, plugin output limits,
process-group cleanup, credentials handling, interpreter dependencies and
restart scheduling behavior. Before web writes or remote exposure: implement
authentication and authorization. History retention and metrics are later work;
no unbounded per-execution log collector is enabled.

## UI direction

Carry forward the useful information architecture of Nagios Modern UI: sidebar,
overview, severity-first problems, host detail, groups and graph drill-down.
Extract shared React components rather than copying PHP pages or retaining CGI
iframes. This initial UI establishes the SMon identity and real readiness only.

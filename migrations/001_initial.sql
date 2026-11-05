CREATE TABLE hosts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    address text NOT NULL CHECK (btrim(address) <> ''),
    enabled boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Explicit executable and argv. Nagios command-line import is a future adapter.
CREATE TABLE commands (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    executable text NOT NULL CHECK (btrim(executable) <> ''),
    arguments text[] NOT NULL DEFAULT '{}',
    timeout_seconds integer NOT NULL DEFAULT 30 CHECK (timeout_seconds > 0)
);

CREATE TABLE services (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    host_id bigint NOT NULL REFERENCES hosts(id),
    command_id bigint NOT NULL REFERENCES commands(id),
    description text NOT NULL CHECK (btrim(description) <> ''),
    arguments text[] NOT NULL DEFAULT '{}',
    check_interval_seconds integer NOT NULL DEFAULT 300 CHECK (check_interval_seconds > 0),
    retry_interval_seconds integer NOT NULL DEFAULT 60 CHECK (retry_interval_seconds > 0),
    max_check_attempts integer NOT NULL DEFAULT 3 CHECK (max_check_attempts > 0),
    enabled boolean NOT NULL DEFAULT true,
    UNIQUE (host_id, description)
);

-- No row means pending: a service is never implicitly OK before its first check.
CREATE TABLE service_states (
    service_id bigint PRIMARY KEY REFERENCES services(id),
    state smallint NOT NULL CHECK (state BETWEEN 0 AND 3),
    state_type text NOT NULL CHECK (state_type IN ('SOFT', 'HARD')),
    attempt integer NOT NULL CHECK (attempt > 0),
    output text NOT NULL,
    long_output text NOT NULL DEFAULT '',
    perfdata text NOT NULL DEFAULT '',
    checked_at timestamptz NOT NULL,
    state_changed_at timestamptz NOT NULL
);

CREATE TABLE state_events (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_id bigint NOT NULL REFERENCES services(id),
    previous_state smallint CHECK (previous_state BETWEEN 0 AND 3),
    state smallint NOT NULL CHECK (state BETWEEN 0 AND 3),
    state_type text NOT NULL CHECK (state_type IN ('SOFT', 'HARD')),
    output text NOT NULL,
    occurred_at timestamptz NOT NULL
);
CREATE INDEX state_events_service_time ON state_events (service_id, occurred_at DESC);

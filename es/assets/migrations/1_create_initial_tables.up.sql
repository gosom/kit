CREATE TYPE "command_status" AS ENUM ('pending', 'running', 'finished', 'failure');

CREATE TABLE "commands" (
    id VARCHAR(26) PRIMARY KEY,
    aggregate_id VARCHAR(50) NOT NULL,
    aggregate_hash INTEGER NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status "command_status" DEFAULT NULL
);

CREATE TABLE "aggregate_versions" (
    aggregate_id VARCHAR(50) PRIMARY KEY,
    version INTEGER NOT NULL
);

CREATE TABLE "events" (
    id VARCHAR(26) PRIMARY KEY,
    command_id VARCHAR(26) NOT NULL REFERENCES "commands" (id),
    aggregate_id VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE TABLE "subscriptions" (
    subscription_group VARCHAR(50) PRIMARY KEY,
    last_event_id VARCHAR(26) DEFAULT NULL REFERENCES "events" (id),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);


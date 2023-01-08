package postgres

const (
	saveCommandsStmt = `
	INSERT INTO "commands" 
		(id, aggregate_id, event_type, data, created_at, aggregate_hash)
	VALUES
		%s
	ON CONFLICT DO NOTHING
	RETURNING id`

	getCommandStmt = `
	SELECT
		id, aggregate_id, event_type, data, created_at, aggregate_hash, status
	FROM
		"commands"
	WHERE
		id = $1`

	selectCommandsToProcess = `
	WITH cte AS (
		SELECT 
		id, aggregate_id, event_type, data, created_at, aggregate_hash, 
		status, ROW_NUMBER() 
		OVER (PARTITION BY MOD(aggregate_hash, $1) ORDER BY id ASC) AS rn
		FROM "commands"
		WHERE status IS NULL
	)
	SELECT 
	id, aggregate_id, event_type, data, created_at, 
	aggregate_hash, COALESCE(status::text, ''), rn
	FROM cte
	WHERE rn <= $2
	`

	checkVersionStmt = `
	UPDATE "aggregate_versions"
	SET version = version + $1
	WHERE aggregate_id = $2 AND version = $3
	`

	saveEventsStmt = `
	INSERT INTO "events"
		(id, command_id, aggregate_id, version, event_type, data)
	VALUES
		($1, $2, $3, $4, $5, $6)
	`
	updateCommandStatusStmt = `
	UPDATE "commands"
		SET status = $1
	WHERE id = $2`

	getOrCreateAggregateVersionStmt = `
	WITH cte AS (
		INSERT INTO "aggregate_versions"
			(aggregate_id, version)
		VALUES
			($1, 0)
		ON CONFLICT (aggregate_id) DO NOTHING
		RETURNING aggregate_id, version
	)
	SELECT aggregate_id, version FROM cte
	UNION
	SELECT aggregate_id, version FROM "aggregate_versions" WHERE aggregate_id = $1`

	insertSubStmt = `
	WITH cte AS (
		INSERT INTO "subscriptions"
		(subscription_group)
		VALUES
		($1)
		ON CONFLICT DO NOTHING
		RETURNING subscription_group, last_event_id, updated_at
	)
	SELECT subscription_group, COALESCE(last_event_id, ''), updated_at FROM cte
	UNION
	SELECT subscription_group, COALESCE(last_event_id, ''), updated_at FROM "subscriptions" WHERE subscription_group = $1`

	selectEventsForSubStmt = `
	WITH cte AS (
	SELECT 
		COALESCE(last_event_id, '') AS last_event_id
	FROM "subscriptions"
	WHERE subscription_group = $1
	)
	SELECT id, aggregate_id, event_type, data, created_at, command_id, version
	FROM events
	WHERE 
	id > (SELECT last_event_id FROM cte)
	AND event_type != 'EventError'
	ORDER BY id, version ASC
	LIMIT $2`

	updateSubStmt = `
	UPDATE "subscriptions"
	SET last_event_id = $2, updated_at = (NOW() at time zone 'utc')
	WHERE subscription_group = $1
	RETURNING subscription_group, last_event_id, updated_at`

	loadEventsStmt = `
	SELECT id, aggregate_id, event_type, data, created_at, command_id, version
	FROM events
	WHERE 
	aggregate_id = $1
	AND event_type != 'EventError'
	ORDER BY id, version ASC
	`
)

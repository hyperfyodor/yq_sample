-- name: CreateTask :one
INSERT INTO tasks (type, value)
VALUES ($1, $2) RETURNING id;

-- name: GetTaskState :one
SELECT state
FROM tasks
WHERE id = $1;

-- name: SetStateToProcessing :one
UPDATE tasks
SET state = 'processing', last_update_time = extract(epoch from now())
WHERE id = $1 RETURNING id;

-- name: SetStateToDone :one
UPDATE tasks
SET state = 'done', last_update_time = extract(epoch from now())
WHERE id = $1 RETURNING id;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
-- name: GetTasks :many
SELECT id, description, completed, addedFrom FROM tasks;

-- name: AddTask :exec
INSERT INTO tasks (description, completed, addedFrom) VALUES (?, ?, ?);

-- name: Checked :exec
UPDATE tasks SET completed = NOT completed WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

-- name: FilterCompletedTasks :many
SELECT id, description, completed, addedFrom FROM tasks WHERE completed = TRUE;

-- name: FilterIncompleteTasks :many
SELECT id, description, completed, addedFrom FROM tasks WHERE completed = FALSE;

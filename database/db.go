// todoapp/database/db.go

package database

import (
	"context"
	"database/sql"
	"fmt"
	"todoapp/models"
)

func GetAllTodos(ctx context.Context, db *sql.DB) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := db.QueryContext(ctx, "SELECT id, title, description, completed, created_at FROM todos ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func GetTodoByID(ctx context.Context, db *sql.DB, id int) (models.Todo, error) {
	var t models.Todo
	row := db.QueryRowContext(ctx, "SELECT id, title, description, completed, created_at FROM todos WHERE id = $1", id)
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
	if err != nil {
		return t, err
	}
	return t, nil
}

func GetTodosByCompletionStatus(ctx context.Context, db *sql.DB, completed bool) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := db.QueryContext(ctx, "SELECT id, title, description, completed, created_at FROM todos WHERE completed = $1 ORDER BY created_at ASC", completed)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func UpdateTodoByID(ctx context.Context, db *sql.DB, id int, title, description string, completed bool) error {
	_, err := db.ExecContext(ctx, "UPDATE todos SET title = $1, description = $2, completed = $3 WHERE id = $4", title, description, completed, id)
	return err
}

func CreateTodo(ctx context.Context, db *sql.DB, title, description string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO todos (title, description, completed, created_at) VALUES ($1, $2, $3, NOW())", title, description, false)
	return err
}

func DeleteTodoByID(ctx context.Context, db *sql.DB, id int) error {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("todo c id %d не найден", id)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}

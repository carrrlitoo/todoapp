// repository/repository.go

package repository

import (
	"context"
	"todoapp/models"
)

type TodoRepository interface {
	GetAllTodos(ctx context.Context) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id int) (models.Todo, error)
	UpdateTodoByID(ctx context.Context, id int, title, description string, completed bool) error
	CreateTodo(ctx context.Context, title, description string) error
	DeleteTodoByID(ctx context.Context, id int) error
	GetTodosByCompletionStatus(ctx context.Context, completed bool) ([]models.Todo, error)
}

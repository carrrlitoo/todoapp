// todoapp/handlers/todos.go

package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todoapp/models"
	"todoapp/repository"
	"todoapp/validation"

	"github.com/go-chi/chi/v5"
)

func GetTodos(todoRepo repository.TodoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		completedParam := r.URL.Query().Get("completed")

		var todos []models.Todo
		var err error

		if completedParam == "true" {
			todos, err = todoRepo.GetTodosByCompletionStatus(ctx, true)
			if err != nil {
				if err == context.DeadlineExceeded {
					http.Error(w, "Request timed out", http.StatusGatewayTimeout)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if completedParam == "false" {
			todos, err = todoRepo.GetTodosByCompletionStatus(ctx, false)
			if err != nil {
				if err == context.DeadlineExceeded {
					http.Error(w, "Request timed out", http.StatusGatewayTimeout)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			todos, err = todoRepo.GetAllTodos(ctx)
			if err != nil {
				if err == context.DeadlineExceeded {
					http.Error(w, "Request timed out", http.StatusGatewayTimeout)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

func GetTodoByID(todoRepo repository.TodoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		todo, err := todoRepo.GetTodoByID(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Todo not found", http.StatusNotFound)
				return
			} else if err == context.DeadlineExceeded {
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
}

func CreateTodo(todoRepo repository.TodoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var todo models.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = validation.IsValidTodoTitle(todo.Title)
		if err != nil {
			if ve, ok := err.(*validation.ValidationError); ok {
				http.Error(w, ve.Message, ve.Code)
				return
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		err = todoRepo.CreateTodo(ctx, todo.Title, todo.Description)
		if err != nil {
			if err == context.DeadlineExceeded {
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func UpdateTodoByID(todoRepo repository.TodoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var todo models.Todo
		err = json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = validation.IsValidTodoTitle(todo.Title)
		fmt.Println("Validation error:", err)
		if HandleValidationError(w, err) {
			return
		}

		err = todoRepo.UpdateTodoByID(ctx, id, todo.Title, todo.Description, todo.Completed)
		if err != nil {
			if err == context.DeadlineExceeded {
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func DeleteTodoByID(todoRepo repository.TodoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = todoRepo.DeleteTodoByID(ctx, id)
		if err != nil {
			if err.Error() == fmt.Sprintf("todo c id %d не найден", id) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if err == context.DeadlineExceeded {
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

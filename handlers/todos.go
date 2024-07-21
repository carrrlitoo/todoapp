// todoapp/handlers/todos.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoapp/database"
	"todoapp/models"
	"todoapp/validation"

	"github.com/go-chi/chi/v5"
)

func GetTodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		completedParam := r.URL.Query().Get("completed")

		var todos []models.Todo
		var err error

		if completedParam == "true" {
			todos, err = database.GetTodosByCompletionStatus(db, true)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if completedParam == "false" {
			todos, err = database.GetTodosByCompletionStatus(db, false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			todos, err = database.GetAllTodos(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

func GetTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		todo, err := database.GetTodoByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Todo not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
}

func CreateTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = database.CreateTodo(db, todo.Title, todo.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func UpdateTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = database.UpdateTodoByID(db, id, todo.Title, todo.Description, todo.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func DeleteTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = database.DeleteTodoByID(db, id)
		if err != nil {
			if err.Error() == fmt.Sprintf("todo c id %d не найден", id) {
				http.Error(w, err.Error(), http.StatusNotFound)
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

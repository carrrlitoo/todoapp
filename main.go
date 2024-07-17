// todoapp/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"todoapp/config"
	"todoapp/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func main() {
	connStr := config.GetDBConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Ошибка при проверке подключения к базе данных:", err)
		return
	}
	fmt.Println("Успешное подключение к базе данных! Successfully connected to the database!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/todos", handlers.GetTodos(db))
	r.Post("/todos", handlers.CreateTodo(db))

	r.Get("/todos/{id}", handlers.GetTodoByID(db))
	r.Put("/todos/{id}", handlers.UpdateTodoByID(db))
	r.Delete("/todos/{id}", handlers.DeleteTodoByID(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}

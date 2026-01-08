// todoapp/main.go

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	r.Use(middleware.Logger)

	r.Get("/todos", handlers.GetTodos(db))
	r.Post("/todos", handlers.CreateTodo(db))

	r.Get("/todos/{id}", handlers.GetTodoByID(db))
	r.Put("/todos/{id}", handlers.UpdateTodoByID(db))
	r.Delete("/todos/{id}", handlers.DeleteTodoByID(db))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	select {
	case <-quit:
		fmt.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %s", err)
		}
		fmt.Println("Server shutdown complete.")
	}
}

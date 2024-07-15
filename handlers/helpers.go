// todoapp/handlers/todos.go

package handlers

import (
	"net/http"
	"todoapp/validation"
)

func HandleValidationError(w http.ResponseWriter, err error) bool {
	if err != nil {
		if validErr, ok := err.(*validation.ValidationError); ok {
			http.Error(w, validErr.Message, validErr.Code)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return true
	}
	return false
}

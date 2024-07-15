// todoapp/validation/todoValidation.go

package validation

import (
	"net/http"
	"strings"
)

type ValidationError struct {
	Code    int
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func IsValidTodoTitle(title string) error {
	if strings.TrimSpace(title) == "" || len(title) == 0 {
		return &ValidationError{
			Code:    http.StatusBadRequest,
			Message: "Title cannot be empty",
		}
	} else if len(title) > 200 {
		return &ValidationError{
			Code:    http.StatusBadRequest,
			Message: "Title cannot exceed 200 characters",
		}
	}
	return nil
}

package model

type ValidationError struct {
	message string
}

// Создание новой ошибки валидации
func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

// Возврат текста ошибки
func (e *ValidationError) Error() string {
	return e.message
}

// Проверка на ошибку валидации
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
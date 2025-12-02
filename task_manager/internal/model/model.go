// internal/model/model.go
package model

// Model - интерфейс для всех моделей проекта
type Model interface {
    GetID() int
    GetType() string
}
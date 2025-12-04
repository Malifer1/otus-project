// internal/model/model.go
package model

// Model - интерфейс для всех моделей
type Model interface {
    GetID() int
    GetType() string
}
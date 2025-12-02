package repository

import "task-manager/internal/model"

// Storage хранит разные типы моделей в отдельных слайсах
type Storage struct {
    Tasks []model.Model
    Notes []model.Model
}

// Генерация нового хранилища
func NewStorage() *Storage {
    return &Storage{
        Tasks: make([]model.Model, 0),
        Notes: make([]model.Model, 0),
    }
}

// Добавление модели в соответствующий слайс на основе её типа
func (s *Storage) AddModel(m model.Model) {
    switch m.GetType() {
    case "task":
        s.Tasks = append(s.Tasks, m)
    case "note":
        s.Notes = append(s.Notes, m)
    }
}

// Подсчет числа моделей каждого типа
func (s *Storage) Count() (int, int) {
    return len(s.Tasks), len(s.Notes)
}
package repository

import (
    "task-manager/internal/model"
)

// Storage - хранилище моделей
type Storage struct {
    tasks []*model.Task  // Конкретный тип для задач
    notes []*model.Note  // Конкретный тип для заметок
}

// NewStorage создаёт новое хранилище
func NewStorage() *Storage {
    return &Storage{
        tasks: make([]*model.Task, 0),
        notes: make([]*model.Note, 0),
    }
}

// AddModel добавляет модель в соответствующий слайс с использованием type switch
func (s *Storage) AddModel(m interface{}) error {
    switch v := m.(type) {
    case *model.Task:
        s.tasks = append(s.tasks, v)
        return nil
    case *model.Note:
        s.notes = append(s.notes, v)
        return nil
    default:
        return model.NewValidationError("unknown model type")
    }
}

// GetTasks возвращает все задачи
func (s *Storage) GetTasks() []*model.Task {
    return s.tasks
}

// GetNotes возвращает все заметки
func (s *Storage) GetNotes() []*model.Note {
    return s.notes
}

// Count возвращает количество моделей каждого типа
func (s *Storage) Count() (int, int) {
    return len(s.tasks), len(s.notes)
}
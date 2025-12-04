// internal/repository/repository.go
package repository

import (
    "sync"
    "task-manager/internal/model"
)

// Storage - потокобезопасное хранилище
type Storage struct {
    tasks []model.Model
    notes []model.Model
    mu    sync.RWMutex
}

// NewStorage создаёт новое хранилище
func NewStorage() *Storage {
    return &Storage{
        tasks: make([]model.Model, 0),
        notes: make([]model.Model, 0),
    }
}

// AddModel добавляет модель в соответствующий слайс
func (s *Storage) AddModel(m model.Model) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    switch m.GetType() {
    case "task":
        s.tasks = append(s.tasks, m)
        return nil
    case "note":
        s.notes = append(s.notes, m)
        return nil
    default:
        return model.NewValidationError("unknown model type")
    }
}

// GetTasks возвращает все задачи
func (s *Storage) GetTasks() []model.Model {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    tasks := make([]model.Model, len(s.tasks))
    copy(tasks, s.tasks)
    return tasks
}

// GetNotes возвращает все заметки
func (s *Storage) GetNotes() []model.Model {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    notes := make([]model.Model, len(s.notes))
    copy(notes, s.notes)
    return notes
}

// Count возвращает количество моделей каждого типа
func (s *Storage) Count() (int, int) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.tasks), len(s.notes)
}

// GetNewTasks возвращает задачи, добавленные после определённого индекса
func (s *Storage) GetNewTasks(lastIndex int) []model.Model {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if lastIndex >= len(s.tasks) {
        return []model.Model{}
    }
    
    newTasks := s.tasks[lastIndex:]
    result := make([]model.Model, len(newTasks))
    copy(result, newTasks)
    return result
}

// GetNewNotes возвращает заметки, добавленные после определённого индекса
func (s *Storage) GetNewNotes(lastIndex int) []model.Model {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if lastIndex >= len(s.notes) {
        return []model.Model{}
    }
    
    newNotes := s.notes[lastIndex:]
    result := make([]model.Model, len(newNotes))
    copy(result, newNotes)
    return result
}
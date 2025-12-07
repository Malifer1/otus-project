package repository

import (
    "sync"
    "task-manager/internal/model"
)

// Storage - потокобезопасное хранилище
type Storage struct {
    tasks []*model.Task
    notes []*model.Note
    mu    sync.RWMutex
}

// NewStorage создаёт новое хранилище
func NewStorage() *Storage {
    return &Storage{
        tasks: make([]*model.Task, 0),
        notes: make([]*model.Note, 0),
    }
}

// AddModel добавляет модель в соответствующий слайс
func (s *Storage) AddModel(m interface{}) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
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

// GetTasks возвращает копию слайса с задачами
func (s *Storage) GetTasks() []*model.Task {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Возвращаем копию для безопасности
    tasks := make([]*model.Task, len(s.tasks))
    copy(tasks, s.tasks)
    return tasks
}

// GetNotes возвращает копию слайса с заметками
func (s *Storage) GetNotes() []*model.Note {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    notes := make([]*model.Note, len(s.notes))
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
func (s *Storage) GetNewTasks(lastIndex int) []*model.Task {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if lastIndex >= len(s.tasks) {
        return []*model.Task{}
    }
    
    newTasks := s.tasks[lastIndex:]
    result := make([]*model.Task, len(newTasks))
    copy(result, newTasks)
    return result
}

// GetNewNotes возвращает заметки, добавленные после определённого индекса
func (s *Storage) GetNewNotes(lastIndex int) []*model.Note {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if lastIndex >= len(s.notes) {
        return []*model.Note{}
    }
    
    newNotes := s.notes[lastIndex:]
    result := make([]*model.Note, len(newNotes))
    copy(result, newNotes)
    return result
}

// Cleanup освобождает ресурсы при завершении
func (s *Storage) Cleanup() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Очищаем слайсы (помогает сборщику мусора)
    s.tasks = nil
    s.notes = nil
}
// internal/service/logger.go
package service

import (
    "log"
    "task-manager/internal/repository"
    "time"
)

// Logger отслеживает изменения в хранилище и логирует новые элементы
func Logger(storage *repository.Storage, interval time.Duration, stopChan <-chan struct{}, done chan<- bool) {
    defer func() {
        done <- true
    }()
    
    log.Println("Логер: запущен")
    
    var lastTaskIndex, lastNoteIndex int
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Проверяем новые задачи
            newTasks := storage.GetNewTasks(lastTaskIndex)
            if len(newTasks) > 0 {
                log.Printf("Логер: найдено %d новых задач\n", len(newTasks))
                for _, task := range newTasks {
                    if t, ok := task.(interface {
                        GetTitle() string
                        GetID() int
                        GetStatus() string
                        GetPriority() string
                    }); ok {
                        log.Printf("  - Задача: %s (ID: %d, Статус: %s, Приоритет: %s)",
                            t.GetTitle(), t.GetID(), t.GetStatus(), t.GetPriority())
                    }
                }
                lastTaskIndex += len(newTasks)
            }
            
            // Проверяем новые заметки
            newNotes := storage.GetNewNotes(lastNoteIndex)
            if len(newNotes) > 0 {
                log.Printf("Логер: найдено %d новых заметок\n", len(newNotes))
                for _, note := range newNotes {
                    if n, ok := note.(interface {
                        GetTitle() string
                        GetID() int
                        GetCategory() string
                    }); ok {
                        log.Printf("  - Заметка: %s (ID: %d, Категория: %s)",
                            n.GetTitle(), n.GetID(), n.GetCategory())
                    }
                }
                lastNoteIndex += len(newNotes)
            }
            
            // Если не найдено новых элементов
            if len(newTasks) == 0 && len(newNotes) == 0 {
                log.Println("Логер: новых элементов не найдено")
            }
            
        case <-stopChan:
            log.Println("Логер: получен сигнал остановки")
            return
        }
    }
}
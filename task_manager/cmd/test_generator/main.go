// cmd/test_generator/main.go
package main

import (
    "fmt"
    "task-manager/internal/repository"
    "task-manager/internal/service"
    "time"
)

func main() {
    fmt.Println("Testing model generator!!!\n")
    
    // Создаём хранилище
    storage := repository.NewStorage()
    
    fmt.Println("Enabling generator...")
    
    go func() {
        service.GenerateModels(storage, 10)
        fmt.Println("Генерация завершена!")
    }()
    
	// Таймаут генератор
    time.Sleep(2 * time.Second)
    
    // Получение результатов
    taskCount, noteCount := storage.Count()
    
    fmt.Println("\nResults:")
    fmt.Printf("Создано задач: %d\n", taskCount)
    fmt.Printf("Создано заметок: %d\n", noteCount)
    
    // Выводим все задачи
    fmt.Println("\nСозданные задачи:")
    for i, model := range storage.Tasks {
        if task, ok := model.(interface{ GetTitle() string }); ok {
            fmt.Printf("%d. %s (ID: %d)\n", i+1, task.GetTitle(), model.GetID())
        }
    }
    
    // Выводим все заметки
    fmt.Println("\nСозданные заметки:")
    for i, model := range storage.Notes {
        if note, ok := model.(interface{ GetTitle() string }); ok {
            fmt.Printf("%d. %s (ID: %d)\n", i+1, note.GetTitle(), model.GetID())
        }
    }
    
    fmt.Println("\nla fin")
}
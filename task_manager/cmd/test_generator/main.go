package main

import (
    "fmt"
    "task-manager/internal/repository"
    "task-manager/internal/service"
    "time"
)

func main() {
    fmt.Println("=== Тестирование генератора моделей ===\n")
    
    storage := repository.NewStorage()
    
    fmt.Println("Запускаем генератор моделей...")
    
    // Запускаем генератор в горутине
    go func() {
        if err := service.GenerateModels(storage, 10); err != nil {
            fmt.Printf("Ошибка генерации: %v\n", err)
        } else {
            fmt.Println("Генерация завершена!")
        }
    }()
    
    // Ждём, чтобы генератор успел поработать
    time.Sleep(2 * time.Second)
    
    // Получаем результаты
    taskCount, noteCount := storage.Count()
    
    fmt.Println("\n=== Результаты ===")
    fmt.Printf("Создано задач: %d\n", taskCount)
    fmt.Printf("Создано заметок: %d\n", noteCount)
    
    // Выводим все задачи
    fmt.Println("\n=== Созданные задачи ===")
    for i, task := range storage.GetTasks() {
        fmt.Printf("%d. %s (ID: %d, Статус: %s)\n", 
            i+1, task.GetTitle(), task.GetID(), task.GetStatus())
    }
    
    // Выводим все заметки
    fmt.Println("\n=== Созданные заметки ===")
    for i, note := range storage.GetNotes() {
        fmt.Printf("%d. %s (ID: %d, Создана: %v)\n", 
            i+1, note.GetTitle(), note.GetID(), 
            note.GetCreatedAt().Format("15:04:05"))
    }
    
    // Тестируем добавление неизвестного типа
    fmt.Println("\n=== Тест обработки ошибок ===")
    if err := storage.AddModel("неизвестный тип"); err != nil {
        fmt.Printf("Ожидаемая ошибка: %v\n", err)
    }
    
    fmt.Println("\n=== Тест завершён ===")
}
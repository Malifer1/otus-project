// cmd/concurrent/main.go
package main

import (
    "fmt"
    "log"
    "task-manager/internal/model"
    "task-manager/internal/repository"
    "task-manager/internal/service"
    "time"
)

func main() {
    fmt.Println("=== Многопоточная система Task Manager ===\n")
    
    // Инициализация
    storage := repository.NewStorage()
    modelChan := make(chan model.Model, 10) // Канал для интерфейса Model
    
    // Каналы для управления горутинами
    doneGenerator := make(chan bool, 1)
    doneReceiver := make(chan bool, 1)
    doneLogger := make(chan bool, 1)
    stopLogger := make(chan struct{})
    
    // Запуск логера (работает каждые 200 мс)
    log.Println("Запуск логера...")
    go service.Logger(storage, 200*time.Millisecond, stopLogger, doneLogger)
    
    // Даём логеру время на запуск
    time.Sleep(50 * time.Millisecond)
    
    // Запуск приёмника
    log.Println("Запуск приёмника...")
    go service.Receiver(modelChan, storage, doneReceiver)
    
    // Даём приёмнику время на запуск
    time.Sleep(50 * time.Millisecond)
    
    // Запуск генератора (создаст 10 моделей)
    log.Println("Запуск генератора...")
    go service.GenerateModels(modelChan, 10, doneGenerator)
    
    // Ожидаем завершения генератора
    <-doneGenerator
    fmt.Println("\nГенератор завершил работу")
    
    // Закрываем канал, чтобы приёмник знал, что данные закончились
    close(modelChan)
    
    // Ожидаем завершения приёмника
    <-doneReceiver
    fmt.Println("Приёмник завершил работу")
    
    // Даём логеру время на финальную проверку
    time.Sleep(500 * time.Millisecond)
    
    // Останавливаем логер
    close(stopLogger)
    
    // Ожидаем завершения логера
    <-doneLogger
    fmt.Println("Логер завершил работу")
    
    // Финальная статистика
    taskCount, noteCount := storage.Count()
    fmt.Printf("\n=== ФИНАЛЬНАЯ СТАТИСТИКА ===\n")
    fmt.Printf("Всего задач: %d\n", taskCount)
    fmt.Printf("Всего заметок: %d\n", noteCount)
    fmt.Printf("Всего моделей: %d\n", taskCount+noteCount)
    
    fmt.Println("\nПрограмма завершена.")
}
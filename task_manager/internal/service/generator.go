// internal/service/generator.go
package service

import (
    "fmt"
    "math/rand"
    "task-manager/internal/model"
    "task-manager/internal/repository"
    "time"
)

// GenerateModels создаёт разные модели и отправляет их в канал
func GenerateModels(modelChan chan<- model.Model, count int, done chan<- bool) {
    defer func() {
        done <- true
    }()
    
    fmt.Println("Генератор: запущен")
    
    rand.Seed(time.Now().UnixNano())
    
    for i := 0; i < count; i++ {
        time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
        
        if i%2 == 0 {
            // Создаём задачу
            dueDate := time.Now().Add(time.Duration(rand.Intn(7)+1) * 24 * time.Hour)
            task, err := model.NewTask(
                fmt.Sprintf("Задача %d", i+1),
                fmt.Sprintf("Описание задачи %d", i+1),
                randomPriority(),
                &dueDate,
            )
            if err != nil {
                fmt.Printf("Генератор: ошибка создания задачи: %v\n", err)
                continue
            }
            task.SetID(i + 1)
            task.SetStatus(randomStatus())
            
            fmt.Printf("Генератор: создана задача '%s'\n", task.GetTitle())
            modelChan <- task
            
        } else {
            // Создаём заметку
            note := model.NewNote(
                fmt.Sprintf("Заметка %d", i+1),
                fmt.Sprintf("Содержимое заметки %d", i+1),
                randomCategory(),
            )
            note.SetID(i + 1)
            
            fmt.Printf("Генератор: создана заметка '%s'\n", note.GetTitle())
            modelChan <- note
        }
    }
    
    fmt.Println("Генератор: завершён")
}

// Receiver получает модели из канала и сохраняет в репозиторий
func Receiver(modelChan <-chan model.Model, storage *repository.Storage, done chan<- bool) {
    defer func() {
        done <- true
    }()
    
    fmt.Println("Приёмник: запущен")
    
    for model := range modelChan {
        if err := storage.AddModel(model); err != nil {
            fmt.Printf("Приёмник: ошибка сохранения модели: %v\n", err)
            continue
        }
        
        // Используем интерфейс Model
        switch model.GetType() {
        case "task":
            if task, ok := model.(interface{ GetTitle() string }); ok {
                fmt.Printf("Приёмник: сохранена задача '%s' (ID: %d)\n", 
                    task.GetTitle(), model.GetID())
            }
        case "note":
            if note, ok := model.(interface{ GetTitle() string }); ok {
                fmt.Printf("Приёмник: сохранена заметка '%s' (ID: %d)\n", 
                    note.GetTitle(), model.GetID())
            }
        }
    }
    
    fmt.Println("Приёмник: завершён")
}

// Вспомогательные функции остаются без изменений
func randomPriority() model.TaskPriority {
    priorities := []model.TaskPriority{
        model.PriorityLow,
        model.PriorityMedium,
        model.PriorityHigh,
    }
    return priorities[rand.Intn(len(priorities))]
}

func randomStatus() model.TaskStatus {
    statuses := []model.TaskStatus{
        model.StatusTodo,
        model.StatusInProgress,
        model.StatusDone,
    }
    return statuses[rand.Intn(len(statuses))]
}

func randomCategory() model.NoteCategory {
    categories := []model.NoteCategory{
        model.CategoryPersonal,
        model.CategoryWork,
        model.CategoryIdea,
    }
    return categories[rand.Intn(len(categories))]
}
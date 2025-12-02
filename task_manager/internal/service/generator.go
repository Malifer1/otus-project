// internal/service/generator.go
package service

import (
    "task-manager/internal/model"
    "task-manager/internal/repository"
    "time"
)

// GenerateModels создаёт разные модели и передаёт их в репозиторий
// Вызывается из main по интервалу
func GenerateModels(repo *repository.Storage, count int) {
    for i := 0; i < count; i++ {
        // По очереди создаем задачу и заметку
        if i%2 == 0 {
            // Создание задачи
            dueDate := time.Now().Add(24 * time.Hour)
            task, _ := model.NewTask(
                "Задача из генератора",
                "Описание сгенерированной задачи",
                model.PriorityMedium,
                &dueDate,
            )
            task.SetID(i + 1)
            repo.AddModel(task)
        } else {
            // Создание заметки
            note := model.NewNote("Заметка из генератора")
            note.SetID(i + 1)
            repo.AddModel(note)
        }
        
        // Таймаут по генерации
        time.Sleep(100 * time.Millisecond)
    }
}
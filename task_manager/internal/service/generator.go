package service

import (
    "task-manager/internal/model"
    "task-manager/internal/repository"
    "time"
)

// GenerateModels создаёт разные модели и передаёт их в репозиторий
func GenerateModels(storage *repository.Storage, count int) error {
    for i := 0; i < count; i++ {
        if i%2 == 0 {
            // Создаём задачу
            dueDate := time.Now().Add(24 * time.Hour)
            task, err := model.NewTask(
                "Задача из генератора",
                "Описание сгенерированной задачи",
                model.PriorityMedium,
                &dueDate,
            )
            if err != nil {
                return err
            }
            task.SetID(len(storage.GetTasks()) + 1)
            
            if err := storage.AddModel(task); err != nil {
                return err
            }
        } else {
            // Создаём заметку
            note := model.NewNote("Заметка из генератора")
            note.SetID(len(storage.GetNotes()) + 1)
            
            if err := storage.AddModel(note); err != nil {
                return err
            }
        }      
        time.Sleep(100 * time.Millisecond)
    }
    return nil
}
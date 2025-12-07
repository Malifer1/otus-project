package service

import (
	"context"
	"fmt"
	"math/rand"
	"task-manager/internal/model"
	"task-manager/internal/repository"
	"time"
)

// GenerateModels создаёт разные модели и отправляет их в канал
// Завершается при отмене контекста или после создания всех моделей
func GenerateModels(ctx context.Context, modelChan chan<- interface{}, count int) {
	fmt.Println("Генератор: запущен")
	
	for i := 0; i < count; i++ {
		// Проверяем, не отменен ли контекст
		select {
		case <-ctx.Done():
			fmt.Println("Генератор: получен сигнал отмены")
			return
		default:
			// Продолжаем работу
		}
		
		// Случайная задержка между созданием моделей
		delay := time.Duration(rand.Intn(300)+100) * time.Millisecond
		
		select {
		case <-ctx.Done():
			fmt.Println("Генератор: получен сигнал отмены во время задержки")
			return
		case <-time.After(delay):
			// Чередуем создание задач и заметок
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
				
				// Отправляем задачу в канал с проверкой контекста
				select {
				case <-ctx.Done():
					fmt.Println("Генератор: получен сигнал отмены при отправке задачи")
					return
				case modelChan <- task:
					fmt.Printf("Генератор: создана задача '%s'\n", task.GetTitle())
				}
			} else {
				// Создаём заметку
				note := model.NewNote(
					fmt.Sprintf("Заметка %d", i+1),
					fmt.Sprintf("Содержимое заметки %d", i+1),
					randomCategory(),
				)
				note.SetID(i + 1)
				
				// Отправляем заметку в канал с проверкой контекста
				select {
				case <-ctx.Done():
					fmt.Println("Генератор: получен сигнал отмены при отправке заметки")
					return
				case modelChan <- note:
					fmt.Printf("Генератор: создана заметка '%s'\n", note.GetTitle())
				}
			}
		}
	}
	
	fmt.Printf("Генератор: успешно создано %d моделей\n", count)
}

// Receiver получает модели из канала и сохраняет в репозиторий
// Завершается при отмене контекста или закрытии канала
func Receiver(ctx context.Context, modelChan <-chan interface{}, storage *repository.Storage) {
	fmt.Println("Приёмник: запущен")
	
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Приёмник: получен сигнал отмены")
			return
		case model, ok := <-modelChan:
			if !ok {
				fmt.Println("Приёмник: канал закрыт")
				return
			}
			
			if err := storage.AddModel(model); err != nil {
				fmt.Printf("Приёмник: ошибка сохранения модели: %v\n", err)
				continue
			}
			
			switch v := model.(type) {
			case interface{ GetTitle() string; GetID() int }:
				fmt.Printf("Приёмник: сохранена модель '%s' (ID: %d)\n", 
					v.GetTitle(), v.GetID())
			}
		}
	}
}

// Вспомогательные функции для генерации случайных данных
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
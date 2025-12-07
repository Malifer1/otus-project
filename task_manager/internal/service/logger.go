package service

import (
	"context"
	"log"
	"task-manager/internal/repository"
	"time"
)

// Logger отслеживает изменения в хранилище и логирует новые элементы
// Завершается при отмене контекста
func Logger(ctx context.Context, storage *repository.Storage, interval time.Duration) {
	log.Println("Логер: запущен")
	
	var lastTaskIndex, lastNoteIndex int
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Логер: получен сигнал отмены")
			// Финальная проверка перед завершением
			logFinalChanges(storage, lastTaskIndex, lastNoteIndex)
			return
		case <-ticker.C:
			logChanges(storage, &lastTaskIndex, &lastNoteIndex)
		}
	}
}

// logChanges проверяет и логирует изменения с момента последней проверки
func logChanges(storage *repository.Storage, lastTaskIndex, lastNoteIndex *int) {
	// Проверяем новые задачи
	newTasks := storage.GetNewTasks(*lastTaskIndex)
	if len(newTasks) > 0 {
		log.Printf("Логер: найдено %d новых задач\n", len(newTasks))
		for _, task := range newTasks {
			log.Printf("  - Задача: %s (ID: %d, Статус: %v, Приоритет: %v)",
				task.GetTitle(), task.GetID(), task.GetStatus(), task.GetPriority())
		}
		*lastTaskIndex += len(newTasks)
	}
	
	// Проверяем новые заметки
	newNotes := storage.GetNewNotes(*lastNoteIndex)
	if len(newNotes) > 0 {
		log.Printf("Логер: найдено %d новых заметок\n", len(newNotes))
		for _, note := range newNotes {
			log.Printf("  - Заметка: %s (ID: %d, Категория: %v)",
				note.GetTitle(), note.GetID(), note.GetCategory())
		}
		*lastNoteIndex += len(newNotes)
	}
	
	// Если не найдено новых элементов
	if len(newTasks) == 0 && len(newNotes) == 0 {
		log.Println("Логер: новых элементов не найдено")
	}
}

// logFinalChanges выполняет финальную проверку изменений перед завершением
func logFinalChanges(storage *repository.Storage, lastTaskIndex, lastNoteIndex int) {
	log.Println("Логер: финальная проверка изменений...")
	
	newTasks := storage.GetNewTasks(lastTaskIndex)
	newNotes := storage.GetNewNotes(lastNoteIndex)
	
	if len(newTasks) > 0 || len(newNotes) > 0 {
		log.Printf("Логер: найдено непротоколированных изменений: %d задач, %d заметок\n", 
			len(newTasks), len(newNotes))
		
		// Логируем непротоколированные задачи
		for _, task := range newTasks {
			log.Printf("  - Непротоколированная задача: %s (ID: %d)",
				task.GetTitle(), task.GetID())
		}
		
		// Логируем непротоколированные заметки
		for _, note := range newNotes {
			log.Printf("  - Непротоколированная заметка: %s (ID: %d)",
				note.GetTitle(), note.GetID())
		}
	} else {
		log.Println("Логер: непротоколированных изменений нет")
	}
}
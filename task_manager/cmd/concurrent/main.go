package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"task-manager/internal/repository"
	"task-manager/internal/service"
	"time"
)

func main() {
	fmt.Println("=== Многопоточная система Task Manager с Graceful Shutdown ===\n")

	// Создаем контекст, который отменяется при получении сигнала
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Инициализация компонентов
	storage := repository.NewStorage()
	modelChan := make(chan interface{}, 10)

	// WaitGroup для отслеживания всех горутин
	var wg sync.WaitGroup

	// Запуск логера (работает каждые 200 мс)
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Логер: запущен")
		service.Logger(ctx, storage, 200*time.Millisecond)
		fmt.Println("Логер: завершен")
	}()

	// Даем логеру время на запуск
	time.Sleep(50 * time.Millisecond)

	// Запуск приемника
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Приемник: запущен")
		service.Receiver(ctx, modelChan, storage)
		fmt.Println("Приемник: завершен")
	}()

	// Даем приемнику время на запуск
	time.Sleep(50 * time.Millisecond)

	// Запуск генератора (создаст 15 моделей)
	wg.Add(1)
	go func() {
		defer func() {
			// Закрываем канал при завершении генератора
			close(modelChan)
			fmt.Println("Генератор: канал закрыт")
			wg.Done()
		}()
		fmt.Println("Генератор: запущен")
		service.GenerateModels(ctx, modelChan, 15)
		fmt.Println("Генератор: завершен")
	}()

	// Отдельная горутина для обработки сигналов ОС
	go func() {
		sig := <-sigChan
		fmt.Printf("\n\nПолучен сигнал: %v. Начинаем graceful shutdown...\n", sig)
		cancel() // Отменяем контекст

		// Даем время на завершение (3 секунды)
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer timeoutCancel()

		// Канал для ожидания завершения горутин
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		// Ждем либо завершения всех горутин, либо истечения таймаута
		select {
		case <-done:
			fmt.Println("Все горутины завершены корректно")
		case <-timeoutCtx.Done():
			fmt.Println("Время ожидания истекло, принудительное завершение")
		}

		// Завершаем программу
		os.Exit(0)
	}()

	fmt.Println("\nСистема запущена. Для завершения нажмите Ctrl+C")
	fmt.Println("Ожидаем завершения работы...")

	// Ждем завершения всех горутин
	wg.Wait()

	// Финальная статистика
	taskCount, noteCount := storage.Count()
	fmt.Printf("\n=== ФИНАЛЬНАЯ СТАТИСТИКА ===\n")
	fmt.Printf("Всего задач: %d\n", taskCount)
	fmt.Printf("Всего заметок: %d\n", noteCount)
	fmt.Printf("Всего моделей: %d\n", taskCount+noteCount)
	fmt.Println("\n=== Программа завершена корректно ===")
}
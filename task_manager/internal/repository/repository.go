package repository

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"task-manager/internal/model"
	"time"
)

// Storage - потокобезопасное хранилище с сохранением в файлы
type Storage struct {
	tasks []*model.Task
	notes []*model.Note
	mu    sync.RWMutex
	
	tasksFile string
	notesFile string
}

// NewStorage создаёт новое хранилище с указанием файлов для сохранения
func NewStorage(tasksFile, notesFile string) *Storage {
	storage := &Storage{
		tasks:     make([]*model.Task, 0),
		notes:     make([]*model.Note, 0),
		tasksFile: tasksFile,
		notesFile: notesFile,
	}
	
	// Загружаем данные из файлов при создании
	storage.loadFromFiles()
	
	return storage
}

// AddModel добавляет модель в соответствующий слайс и сохраняет в файл
func (s *Storage) AddModel(m interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	switch v := m.(type) {
	case *model.Task:
		s.tasks = append(s.tasks, v)
		// Сохраняем задачи в файл
		if err := s.saveTasksToFile(); err != nil {
			return fmt.Errorf("ошибка сохранения задач: %w", err)
		}
		return nil
	case *model.Note:
		s.notes = append(s.notes, v)
		// Сохраняем заметки в файл
		if err := s.saveNotesToFile(); err != nil {
			return fmt.Errorf("ошибка сохранения заметок: %w", err)
		}
		return nil
	default:
		return model.NewValidationError("unknown model type")
	}
}

// SaveAll сохраняет все данные во все файлы
func (s *Storage) SaveAll() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if err := s.saveTasksToFile(); err != nil {
		return fmt.Errorf("ошибка сохранения задач: %w", err)
	}
	
	if err := s.saveNotesToFile(); err != nil {
		return fmt.Errorf("ошибка сохранения заметок: %w", err)
	}
	
	return nil
}

// loadFromFiles загружает данные из файлов при старте
func (s *Storage) loadFromFiles() {
	// Загружаем задачи
	if err := s.loadTasksFromFile(); err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
	}
	
	// Загружаем заметки
	if err := s.loadNotesFromFile(); err != nil {
		fmt.Printf("Ошибка загрузки заметок: %v\n", err)
	}
}

// ========== Методы для работы с задачами ==========

// saveTasksToFile сохраняет задачи в файлы (CSV и JSON)
func (s *Storage) saveTasksToFile() error {
	// Сохраняем в CSV
	if err := s.saveTasksToCSV(); err != nil {
		return err
	}
	
	// Сохраняем в JSON
	return s.saveTasksToJSON()
}

// saveTasksToCSV сохраняет задачи в CSV файл
func (s *Storage) saveTasksToCSV() error {
	file, err := os.Create(s.tasksFile + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Записываем заголовки
	headers := []string{"ID", "Title", "Description", "Status", "Priority", "CreatedAt", "UpdatedAt", "DueDate"}
	if err := writer.Write(headers); err != nil {
		return err
	}
	
	// Записываем данные задач
	for _, task := range s.tasks {
		var dueDateStr string
		if dueDate := task.GetDueDate(); dueDate != nil {
			dueDateStr = dueDate.Format(time.RFC3339)
		}
		
		record := []string{
			strconv.Itoa(task.GetID()),
			task.GetTitle(),
			task.GetDescription(),
			string(task.GetStatus()),
			string(task.GetPriority()),
			task.GetCreatedAt().Format(time.RFC3339),
			task.GetUpdatedAt().Format(time.RFC3339),
			dueDateStr,
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	
	return nil
}

// saveTasksToJSON сохраняет задачи в JSON файл
func (s *Storage) saveTasksToJSON() error {
	file, err := os.Create(s.tasksFile + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	// Создаем структуру для JSON сериализации
	type jsonTask struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		Priority    string     `json:"priority"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		DueDate     *time.Time `json:"due_date,omitempty"`
	}
	
	var jsonTasks []jsonTask
	for _, task := range s.tasks {
		jsonTasks = append(jsonTasks, jsonTask{
			ID:          task.GetID(),
			Title:       task.GetTitle(),
			Description: task.GetDescription(),
			Status:      string(task.GetStatus()),
			Priority:    string(task.GetPriority()),
			CreatedAt:   task.GetCreatedAt(),
			UpdatedAt:   task.GetUpdatedAt(),
			DueDate:     task.GetDueDate(),
		})
	}
	
	return encoder.Encode(jsonTasks)
}

// loadTasksFromFile загружает задачи из файлов
func (s *Storage) loadTasksFromFile() error {
	// Сначала пробуем загрузить из JSON
	if err := s.loadTasksFromJSON(); err == nil {
		return nil
	}
	
	// Если не получилось, загружаем из CSV
	return s.loadTasksFromCSV()
}

// loadTasksFromCSV загружает задачи из CSV файла
func (s *Storage) loadTasksFromCSV() error {
	file, err := os.Open(s.tasksFile + ".csv")
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не существует - это нормально при первом запуске
		}
		return err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	
	// Пропускаем заголовок
	if len(records) <= 1 {
		return nil
	}
	
	for _, record := range records[1:] {
		if len(record) < 8 {
			continue
		}
		
		id, _ := strconv.Atoi(record[0])
		title := record[1]
		description := record[2]
		status := model.TaskStatus(record[3])
		priority := model.TaskPriority(record[4])
		
		var dueDate *time.Time
		if record[7] != "" {
			parsedDate, err := time.Parse(time.RFC3339, record[7])
			if err == nil {
				dueDate = &parsedDate
			}
		}
		
		task, err := model.NewTask(title, description, priority, dueDate)
		if err != nil {
			fmt.Printf("Ошибка создания задачи из CSV: %v\n", err)
			continue
		}
		
		task.SetID(id)
		task.SetStatus(status)
		
		// Устанавливаем даты из файла
		if createdAt, err := time.Parse(time.RFC3339, record[5]); err == nil {
			task.SetCreatedAt(createdAt)
		}
		
		if updatedAt, err := time.Parse(time.RFC3339, record[6]); err == nil {
			task.SetUpdatedAt(updatedAt)
		}
		
		s.tasks = append(s.tasks, task)
	}
	
	return nil
}

// loadTasksFromJSON загружает задачи из JSON файла
func (s *Storage) loadTasksFromJSON() error {
	file, err := os.Open(s.tasksFile + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	
	// Структура для десериализации
	type jsonTask struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		Priority    string     `json:"priority"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		DueDate     *time.Time `json:"due_date,omitempty"`
	}
	
	var jsonTasks []jsonTask
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonTasks); err != nil {
		return err
	}
	
	for _, jt := range jsonTasks {
		task, err := model.NewTask(jt.Title, jt.Description, model.TaskPriority(jt.Priority), jt.DueDate)
		if err != nil {
			fmt.Printf("Ошибка создания задачи из JSON: %v\n", err)
			continue
		}
		
		task.SetID(jt.ID)
		task.SetStatus(model.TaskStatus(jt.Status))
		task.SetCreatedAt(jt.CreatedAt)
		task.SetUpdatedAt(jt.UpdatedAt)
		
		s.tasks = append(s.tasks, task)
	}
	
	return nil
}

// ========== Методы для работы с заметками ==========

// saveNotesToFile сохраняет заметки в файлы (CSV и JSON)
func (s *Storage) saveNotesToFile() error {
	// Сохраняем в CSV
	if err := s.saveNotesToCSV(); err != nil {
		return err
	}
	
	// Сохраняем в JSON
	return s.saveNotesToJSON()
}

// saveNotesToCSV сохраняет заметки в CSV файл
func (s *Storage) saveNotesToCSV() error {
	file, err := os.Create(s.notesFile + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Записываем заголовки
	headers := []string{"ID", "Title", "Content", "Category", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(headers); err != nil {
		return err
	}
	
	// Записываем данные заметок
	for _, note := range s.notes {
		record := []string{
			strconv.Itoa(note.GetID()),
			note.GetTitle(),
			note.GetContent(),
			note.GetCategory(),
			note.GetCreatedAt().Format(time.RFC3339),
			note.GetUpdatedAt().Format(time.RFC3339),
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	
	return nil
}

// saveNotesToJSON сохраняет заметки в JSON файл
func (s *Storage) saveNotesToJSON() error {
	file, err := os.Create(s.notesFile + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	// Структура для JSON сериализации
	type jsonNote struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Category  string    `json:"category"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	
	var jsonNotes []jsonNote
	for _, note := range s.notes {
		jsonNotes = append(jsonNotes, jsonNote{
			ID:        note.GetID(),
			Title:     note.GetTitle(),
			Content:   note.GetContent(),
			Category:  note.GetCategory(),
			CreatedAt: note.GetCreatedAt(),
			UpdatedAt: note.GetUpdatedAt(),
		})
	}
	
	return encoder.Encode(jsonNotes)
}

// loadNotesFromFile загружает заметки из файлов
func (s *Storage) loadNotesFromFile() error {
	// Сначала пробуем загрузить из JSON
	if err := s.loadNotesFromJSON(); err == nil {
		return nil
	}
	
	// Если не получилось, загружаем из CSV
	return s.loadNotesFromCSV()
}

// loadNotesFromCSV загружает заметки из CSV файла
func (s *Storage) loadNotesFromCSV() error {
	file, err := os.Open(s.notesFile + ".csv")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	
	if len(records) <= 1 {
		return nil
	}
	
	for _, record := range records[1:] {
		if len(record) < 6 {
			continue
		}
		
		id, _ := strconv.Atoi(record[0])
		title := record[1]
		content := record[2]
		category := model.NoteCategory(record[3])
		
		note := model.NewNote(title, content, category)
		note.SetID(id)
		
		// Устанавливаем даты из файла
		if createdAt, err := time.Parse(time.RFC3339, record[4]); err == nil {
			note.SetCreatedAt(createdAt)
		}
		
		if updatedAt, err := time.Parse(time.RFC3339, record[5]); err == nil {
			note.SetUpdatedAt(updatedAt)
		}
		
		s.notes = append(s.notes, note)
	}
	
	return nil
}

// loadNotesFromJSON загружает заметки из JSON файла
func (s *Storage) loadNotesFromJSON() error {
	file, err := os.Open(s.notesFile + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	
	// Структура для десериализации
	type jsonNote struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Category  string    `json:"category"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	
	var jsonNotes []jsonNote
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonNotes); err != nil {
		return err
	}
	
	for _, jn := range jsonNotes {
		note := model.NewNote(jn.Title, jn.Content, model.NoteCategory(jn.Category))
		note.SetID(jn.ID)
		note.SetCreatedAt(jn.CreatedAt)
		note.SetUpdatedAt(jn.UpdatedAt)
		
		s.notes = append(s.notes, note)
	}
	
	return nil
}

// ========== Существующие методы (без изменений) ==========

// GetTasks возвращает копию слайса с задачами
func (s *Storage) GetTasks() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tasks := make([]*model.Task, len(s.tasks))
	copy(tasks, s.tasks)
	return tasks
}

// GetNotes возвращает копию слайса с заметками
func (s *Storage) GetNotes() []*model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	notes := make([]*model.Note, len(s.notes))
	copy(notes, s.notes)
	return notes
}

// Count возвращает количество моделей каждого типа
func (s *Storage) Count() (int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks), len(s.notes)
}

// GetNewTasks возвращает задачи, добавленные после определённого индекса
func (s *Storage) GetNewTasks(lastIndex int) []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if lastIndex >= len(s.tasks) {
		return []*model.Task{}
	}
	
	newTasks := s.tasks[lastIndex:]
	result := make([]*model.Task, len(newTasks))
	copy(result, newTasks)
	return result
}

// GetNewNotes возвращает заметки, добавленные после определённого индекса
func (s *Storage) GetNewNotes(lastIndex int) []*model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if lastIndex >= len(s.notes) {
		return []*model.Note{}
	}
	
	newNotes := s.notes[lastIndex:]
	result := make([]*model.Note, len(newNotes))
	copy(result, newNotes)
	return result
}

// Cleanup освобождает ресурсы и сохраняет данные перед завершением
func (s *Storage) Cleanup() {
	// Сохраняем все данные перед завершением
	if err := s.SaveAll(); err != nil {
		fmt.Printf("Ошибка сохранения данных при завершении: %v\n", err)
	}
	
	// Очищаем слайсы
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = nil
	s.notes = nil
}
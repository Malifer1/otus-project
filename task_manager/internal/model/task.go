package model

import (
	"time"
)

// Основная структура задачи
type Task struct {
	id          int
	title       string
	description string
	status      TaskStatus
	priority    TaskPriority
	createdAt   time.Time
	updatedAt   time.Time
	dueDate     *time.Time
}

// тип данных TaskStatus овтечает за статус задачи
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in-progress"
	StatusDone       TaskStatus = "done"
)

// тип данных TaskPriority отвечает за приоритет задачи
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// ключевая функция по созданию новой задачи
func NewTask(title, description string, priority TaskPriority, dueDate *time.Time) (*Task, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	
	if err := validatePriority(priority); err != nil {
		return nil, err
	}

	now := time.Now()
	task := &Task{
		title:       title,
		description: description,
		status:      StatusTodo,
		priority:    priority,
		createdAt:   now,
		updatedAt:   now,
		dueDate:     dueDate,
	}

	return task, nil
}

// Геттер для ID
func (t *Task) GetID() int {
	return t.id
}

// Сеттер для ID
func (t *Task) SetID(id int) {
	t.id = id
	t.updatedAt = time.Now()
}

// Геттер для названия задачи
func (t *Task) GetTitle() string {
	return t.title
}

// Сеттер для задачи
func (t *Task) SetTitle(title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil
}

// Геттер для описания задачи
func (t *Task) GetDescription() string {
	return t.description
}

// Сеттер для описания задачи
func (t *Task) SetDescription(description string) {
	t.description = description
	t.updatedAt = time.Now()
}

// Геттер статуса задачи
func (t *Task) GetStatus() TaskStatus {
	return t.status
}

// Сеттер статуса задачи
func (t *Task) SetStatus(status TaskStatus) error {
	if err := validateStatus(status); err != nil {
		return err
	}
	t.status = status
	t.updatedAt = time.Now()
	return nil
}

// Геттер приоритета задачи
func (t *Task) GetPriority() TaskPriority {
	return t.priority
}

// Сеттер приоритета задачи
func (t *Task) SetPriority(priority TaskPriority) error {
	if err := validatePriority(priority); err != nil {
		return err
	}
	t.priority = priority
	t.updatedAt = time.Now()
	return nil
}

// Геттер даты создания задачи
func (t *Task) GetCreatedAt() time.Time {
	return t.createdAt
}

// Геттер даты апдейта задачи
func (t *Task) GetUpdatedAt() time.Time {
	return t.updatedAt
}

// Геттер срока выполнения задачи
func (t *Task) GetDueDate() *time.Time {
	return t.dueDate
}

// Сеттер срока выполнения задачи
func (t *Task) SetDueDate(dueDate *time.Time) {
	t.dueDate = dueDate
	t.updatedAt = time.Now()
}

// Выставление статуса "в процессе"
func (t *Task) MarkInProgress() error {
	return t.SetStatus(StatusInProgress)
}

// Выставление статуса "выполнено"
func (t *Task) MarkDone() error {
	return t.SetStatus(StatusDone)
}

// валидация просроченности задачи
func (t *Task) IsOverdue() bool {
	if t.dueDate == nil {
		return false
	}
	return time.Now().After(*t.dueDate)
}

// Геттер количества дней до дедлайна
// Положительное значение - дни ДО дедлайн, отрицательное - дни ПОСЛЕ дедлайна
func (t *Task) DaysUntilDue() *int {
	if t.dueDate == nil {
		return nil
	}
	
	diff := t.dueDate.Sub(time.Now())
	days := int(diff.Hours() / 24)
	return &days
}

// Валидационные функции
func validateTitle(title string) error {
	if title == "" {
		return NewValidationError("title cannot be empty")
	}
	if len(title) > 64 {
		return NewValidationError("title cannot be longer than 64 characters")
	}
	return nil
}

func validateStatus(status TaskStatus) error {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return nil
	default:
		return NewValidationError("invalid task status")
	}
}

func validatePriority(priority TaskPriority) error {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return NewValidationError("invalid task priority")
	}
}

// Геттер типа модели
func (t *Task) GetType() string {
    return "task"
}
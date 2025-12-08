package model

import (
    "time"
)

// Task представляет собой задачу в системе управления задачами
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

// TaskStatus представляет статус задачи
type TaskStatus string

const (
    StatusTodo       TaskStatus = "todo"
    StatusInProgress TaskStatus = "in-progress"
    StatusDone       TaskStatus = "done"
)

// TaskPriority представляет приоритет задачи
type TaskPriority string

const (
    PriorityLow    TaskPriority = "low"
    PriorityMedium TaskPriority = "medium"
    PriorityHigh   TaskPriority = "high"
)

// NewTask создает новую задачу с валидацией входных данных
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

// GetID возвращает идентификатор задачи
func (t *Task) GetID() int {
    return t.id
}

// SetID устанавливает идентификатор задачи (только для внутреннего использования)
func (t *Task) SetID(id int) {
    t.id = id
    t.updatedAt = time.Now()
}

// GetTitle возвращает заголовок задачи
func (t *Task) GetTitle() string {
    return t.title
}

// SetTitle устанавливает заголовок задачи с валидацией
func (t *Task) SetTitle(title string) error {
    if err := validateTitle(title); err != nil {
        return err
    }
    t.title = title
    t.updatedAt = time.Now()
    return nil
}

// GetDescription возвращает описание задачи
func (t *Task) GetDescription() string {
    return t.description
}

// SetDescription устанавливает описание задачи
func (t *Task) SetDescription(description string) {
    t.description = description
    t.updatedAt = time.Now()
}

// GetStatus возвращает статус задачи
func (t *Task) GetStatus() TaskStatus {
    return t.status
}

// SetStatus устанавливает статус задачи с валидацией
func (t *Task) SetStatus(status TaskStatus) error {
    if err := validateStatus(status); err != nil {
        return err
    }
    t.status = status
    t.updatedAt = time.Now()
    return nil
}

// GetPriority возвращает приоритет задачи
func (t *Task) GetPriority() TaskPriority {
    return t.priority
}

// SetPriority устанавливает приоритет задачи с валидацией
func (t *Task) SetPriority(priority TaskPriority) error {
    if err := validatePriority(priority); err != nil {
        return err
    }
    t.priority = priority
    t.updatedAt = time.Now()
    return nil
}

// GetCreatedAt возвращает дату создания задачи
func (t *Task) GetCreatedAt() time.Time {
    return t.createdAt
}

// GetUpdatedAt возвращает дату последнего обновления задачи
func (t *Task) GetUpdatedAt() time.Time {
    return t.updatedAt
}

// GetDueDate возвращает срок выполнения задачи
func (t *Task) GetDueDate() *time.Time {
    return t.dueDate
}

// SetDueDate устанавливает срок выполнения задачи
func (t *Task) SetDueDate(dueDate *time.Time) {
    t.dueDate = dueDate
    t.updatedAt = time.Now()
}

// MarkInProgress помечает задачу как "в процессе"
func (t *Task) MarkInProgress() error {
    return t.SetStatus(StatusInProgress)
}

// MarkDone помечает задачу как "выполнено"
func (t *Task) MarkDone() error {
    return t.SetStatus(StatusDone)
}

// IsOverdue проверяет, просрочена ли задача
func (t *Task) IsOverdue() bool {
    if t.dueDate == nil {
        return false
    }
    return time.Now().After(*t.dueDate)
}

// DaysUntilDue возвращает количество дней до дедлайна
// Положительное значение - дней до дедлайна, отрицательное - дней просрочки
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
    if len(title) > 100 {
        return NewValidationError("title cannot be longer than 100 characters")
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

func (t *Task) GetType() string {
    return "task"
}

// SetCreatedAt устанавливает дату создания задачи
func (t *Task) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

// SetUpdatedAt устанавливает дату последнего обновления задачи
func (t *Task) SetUpdatedAt(updatedAt time.Time) {
	t.updatedAt = updatedAt
}
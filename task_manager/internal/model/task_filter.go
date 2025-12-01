package model

import "time"

// тип для фильтрации задач
type TaskFilter struct {
	status   *TaskStatus
	priority *TaskPriority
	fromDate *time.Time
	toDate   *time.Time
}

// создание фильтра задач
func NewTaskFilter() *TaskFilter {
	return &TaskFilter{}
}

// выставление фильтра по статусу
func (f *TaskFilter) WithStatus(status TaskStatus) *TaskFilter {
	f.status = &status
	return f
}

// выставление фильтра фильтр по приоритету
func (f *TaskFilter) WithPriority(priority TaskPriority) *TaskFilter {
	f.priority = &priority
	return f
}

// выставление фильтра по датам
func (f *TaskFilter) WithDateRange(from, to *time.Time) *TaskFilter {
	f.fromDate = from
	f.toDate = to
	return f
}

// Геттер фильтра статуса
func (f *TaskFilter) GetStatus() *TaskStatus {
	return f.status
}

// Геттер фильтра приоритета
func (f *TaskFilter) GetPriority() *TaskPriority {
	return f.priority
}

// Геттер стартовой даты фильтра
func (f *TaskFilter) GetFromDate() *time.Time {
	return f.fromDate
}

// Геттер конечной даты фильтра
func (f *TaskFilter) GetToDate() *time.Time {
	return f.toDate
}
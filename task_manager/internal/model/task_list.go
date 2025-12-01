package model

// Список задач с методами по фильтрации
type TaskList struct {
	tasks []*Task
}

// Создание списка задач
func NewTaskList() *TaskList {
	return &TaskList{
		tasks: make([]*Task, 0),
	}
}

// Добавление задачи в список
func (tl *TaskList) Add(task *Task) {
	tl.tasks = append(tl.tasks, task)
}

// Удаление задачи из списка по ID
func (tl *TaskList) Remove(taskID int) bool {
	for i, task := range tl.tasks {
		if task.GetID() == taskID {
			tl.tasks = append(tl.tasks[:i], tl.tasks[i+1:]...)
			return true
		}
	}
	return false
}

// Геттер задачи по ID
func (tl *TaskList) GetByID(taskID int) *Task {
	for _, task := range tl.tasks {
		if task.GetID() == taskID {
			return task
		}
	}
	return nil
}

// Геттер всех задач
func (tl *TaskList) GetAll() []*Task {
	return tl.tasks
}

// Геттер на задачи по фильтру
func (tl *TaskList) Filter(filter *TaskFilter) []*Task {
	if filter == nil {
		return tl.tasks
	}

	var result []*Task
	for _, task := range tl.tasks {
		if matchesFilter(task, filter) {
			result = append(result, task)
		}
	}
	return result
}

// Геттер числа задач
func (tl *TaskList) Count() int {
	return len(tl.tasks)
}

// Геттер числа задач по статусу
func (tl *TaskList) CountByStatus(status TaskStatus) int {
	count := 0
	for _, task := range tl.tasks {
		if task.GetStatus() == status {
			count++
		}
	}
	return count
}

// Валидация соответствия задачи фильтру
func matchesFilter(task *Task, filter *TaskFilter) bool {
	if filter.GetStatus() != nil && task.GetStatus() != *filter.GetStatus() {
		return false
	}
	
	if filter.GetPriority() != nil && task.GetPriority() != *filter.GetPriority() {
		return false
	}
	
	if filter.GetFromDate() != nil && task.GetCreatedAt().Before(*filter.GetFromDate()) {
		return false
	}
	
	if filter.GetToDate() != nil && task.GetCreatedAt().After(*filter.GetToDate()) {
		return false
	}
	
	return true
}
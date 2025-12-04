package model

import "time"

// Note представляет заметку в системе
type Note struct {
    id        int
    title     string
    content   string
    category  NoteCategory
    createdAt time.Time
    updatedAt time.Time
}

// NoteCategory представляет категорию заметки
type NoteCategory string

const (
    CategoryPersonal NoteCategory = "personal"
    CategoryWork     NoteCategory = "work"
    CategoryIdea     NoteCategory = "idea"
)

// NewNote создаёт новую заметку
func NewNote(title, content string, category NoteCategory) *Note {
    now := time.Now()
    return &Note{
        title:     title,
        content:   content,
        category:  category,
        createdAt: now,
        updatedAt: now,
    }
}

// GetID возвращает идентификатор заметки
func (n *Note) GetID() int {
    return n.id
}

// SetID устанавливает идентификатор заметки
func (n *Note) SetID(id int) {
    n.id = id
    n.updatedAt = time.Now()
}

// GetTitle возвращает заголовок заметки
func (n *Note) GetTitle() string {
    return n.title
}

// GetContent возвращает содержимое заметки
func (n *Note) GetContent() string {
    return n.content
}

// GetCategory возвращает категорию заметки
func (n *Note) GetCategory() string {
    return string(n.category)
}

// GetCreatedAt возвращает дату создания заметки
func (n *Note) GetCreatedAt() time.Time {
    return n.createdAt
}

// GetUpdatedAt возвращает дату последнего обновления заметки
func (n *Note) GetUpdatedAt() time.Time {
    return n.updatedAt
}

func (n *Note) GetType() string {
    return "note"
}
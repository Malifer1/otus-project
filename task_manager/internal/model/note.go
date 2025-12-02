// internal/model/note.go
package model

import "time"

// Note - новая и упрощенная структура для заметок
type Note struct {
    id        int
    title     string
    createdAt time.Time
}

// NewNote создаёт новую заметку
func NewNote(title string) *Note {
    return &Note{
        title:     title,
        createdAt: time.Now(),
    }
}

func (n *Note) GetID() int {
    return n.id
}

func (n *Note) SetID(id int) {
    n.id = id
}

func (n *Note) GetType() string {
    return "note"
}

func (n *Note) GetTitle() string {
    return n.title
}

func (n *Note) GetCreatedAt() time.Time {
    return n.createdAt
}
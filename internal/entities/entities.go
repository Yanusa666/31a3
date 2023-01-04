package entities

import "time"

type Author struct {
	Id   uint64 `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name,omitempty"`
}

type Post struct {
	Id        uint64    `json:"id" db:"id"`
	AuthorId  uint64    `json:"author_id" db:"author_id"`
	Title     string    `json:"title" db:"title,omitempty"`
	Content   string    `json:"content" db:"content,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

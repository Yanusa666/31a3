package handlers

import (
	"crud/internal/entities"
	"time"
)

type ErrorResp struct {
	Error string `json:"error"`
}

type AddAuthorReq struct {
	Name string `json:"name"`
}

type ListAuthorsResp struct {
	Authors []entities.Author `json:"authors"`
}

type UpdateAuthorReq struct {
	Name string `json:"name"`
}

type AddPostReq struct {
	AuthorId  uint64    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ListPostsResp struct {
	Posts []entities.Post `json:"posts"`
}

type UpdatePostReq struct {
	AuthorId  uint64    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

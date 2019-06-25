package model

import (
	"errors"
	"time"
)

// Product-Entity which is used to store a product object into the database
type Article struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateArticle struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

func (update *UpdateArticle) IsValid() bool {
	return true // Could add some validation logic here
}

func (update *UpdateArticle) TransformToArticle() (*Article, error) {
	if !update.IsValid() {
		return nil, errors.New("update article object is not valid")
	}
	return &Article{
		Description: update.Description,
		Content:     update.Content,
	}, nil
}

package model

import (
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
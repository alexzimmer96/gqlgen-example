package model

import (
	"errors"
	"time"
)

// Product-Entity which is used to store a product object into the database
type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

//======================================================================================================================

type CreateArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func (creationRequest *CreateArticle) IsValid() bool {
	return true // Could add some validation logic here
}

func (creationRequest *CreateArticle) TransformToArticle() (*Article, error) {
	if !creationRequest.IsValid() {
		return nil, errors.New("article object is not valid")
	}
	return &Article{
		Title:       creationRequest.Title,
		Description: creationRequest.Description,
		Content:     creationRequest.Content,
	}, nil
}

//======================================================================================================================

type UpdateArticle struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
}

func (update *UpdateArticle) IsValid() bool {
	return true // Could add some validation logic here
}

func (update *UpdateArticle) MergeChanges(article *Article) (*Article, error) {
	if !update.IsValid() {
		return nil, errors.New("article object is not valid")
	}
	if update.Title != nil {
		article.Title = *update.Title
	}
	if update.Description != nil {
		article.Title = *update.Description
	}
	if update.Content != nil {
		article.Title = *update.Content
	}
	return article, nil
}

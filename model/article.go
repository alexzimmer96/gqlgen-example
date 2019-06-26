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

// Returns if the UpdateArticle-Object is valid
// Could be improved by returning a list of validation errors instead of just a boolean.
func (creationRequest *CreateArticle) IsValid() bool {
	return true // Could add some validation logic here
}

// Transforms a CreateArticle-Request into an Article-Object.
// Returns the new Article-Object or an error, if the CreateArticle-Request is not valid.
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

// Returns if the UpdateArticle-Object is valid
// Could be improved by returning a list of validation errors instead of just a boolean.
func (update *UpdateArticle) IsValid() bool {
	return true // Could add some validation logic here
}

// Merge changes from the UpdateArticle-Request into an existing Article.
// Returns the modified Article or an error, if the UpdateArticle-Object is not valid.
func (update *UpdateArticle) MergeChanges(article *Article) (*Article, error) {
	if !update.IsValid() {
		return nil, errors.New("article object is not valid")
	}
	article.Title = setStringIfNotNil(article.Title, update.Title)
	article.Description = setStringIfNotNil(article.Description, update.Description)
	article.Content = setStringIfNotNil(article.Content, update.Content)
	return article, nil
}

func setStringIfNotNil(oldValue string, newValue *string) string {
	if newValue != nil {
		return *newValue
	}
	return oldValue
}

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	unchanged = "unchanged"
	changed   = "changed"
)

func TestUpdateArticle_MergeChanges(t *testing.T) {
	article := &Article{
		Title:       unchanged,
		Description: unchanged,
		Content:     unchanged,
	}
	ChangeTitle := &UpdateArticle{
		Title: &changed,
	}
	article, err := ChangeTitle.MergeChanges(article)
	assert.NoError(t, err)
	assert.Equal(t, changed, article.Title)
	assert.Equal(t, unchanged, article.Description)
	assert.Equal(t, unchanged, article.Content)
}

package repository

import (
	"encoding/json"
	"fmt"
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/muesli/cache2go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArticleRepository_GetAll(t *testing.T) {
	db := cache2go.Cache("example")
	repo := NewArticleRepository(db)

	article, err := repo.Save(&model.Article{
		Description: "Some Description",
		Content: "Some Content",
	})
	assert.NoError(t, err)
	assert.NotNil(t, article)

	articles := repo.GetAll()
	marshalled, _ := json.MarshalIndent(articles, "", "    ")
	fmt.Printf(">>>>> Articles stored:\n%s\n", marshalled)
}
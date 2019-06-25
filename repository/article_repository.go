package repository

import (
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/gofrs/uuid"
	"github.com/muesli/cache2go"
	"time"
)

type IArticleRepository interface {
	Save(article *model.Article) (*model.Article, error)
	GetAll() []*model.Article
	Delete(string) (bool, error)
}

type ArticleRepository struct {
	db *cache2go.CacheTable
}

func NewArticleRepository(db *cache2go.CacheTable) *ArticleRepository {
	return &ArticleRepository{db}
}

func (repo *ArticleRepository) Save(article *model.Article) (*model.Article, error) {
	now := time.Now()
	if article.ID == "" {
		article.ID = uuid.Must(uuid.NewV4()).String()
		article.CreatedAt = now
	}
	article.UpdatedAt = now
	repo.db.Add(article.ID, 0, article)
	return article, nil
}

func (repo *ArticleRepository) GetAll() []*model.Article {
	var articles []*model.Article
	repo.db.Foreach(func(id interface{}, item *cache2go.CacheItem) {
		articles = append(articles, item.Data().(*model.Article))
	})
	return articles
}

func (repo *ArticleRepository) Delete(id string) (bool, error) {
	_, err := repo.db.Delete(id)
	if err == cache2go.ErrKeyNotFound {
		return false, err
	}
	return true, nil
}
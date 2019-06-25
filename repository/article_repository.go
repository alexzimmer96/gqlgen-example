package repository

import (
	"fmt"
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/gofrs/uuid"
	"github.com/muesli/cache2go"
	"github.com/pkg/errors"
	"time"
)

type IArticleRepository interface {
	Save(article *model.Article) (*model.Article, error)
	GetAll() ([]*model.Article, error)
	GetSingle(string) (*model.Article, error)
	Delete(string) (bool, error)
	GetCreationStream() chan *model.Article
}

type ArticleRepository struct {
	db             *cache2go.CacheTable
	creationStream chan *model.Article
}

func NewArticleRepository(db *cache2go.CacheTable) *ArticleRepository {
	return &ArticleRepository{
		db:             db,
		creationStream: make(chan *model.Article),
	}
}

func (repo *ArticleRepository) Save(article *model.Article) (*model.Article, error) {
	created := false
	now := time.Now()
	if article.ID == "" {
		article.ID = uuid.Must(uuid.NewV4()).String()
		article.CreatedAt = now
		created = true
	}
	article.UpdatedAt = now
	repo.db.Add(article.ID, 0, article)
	if created {
		repo.creationStream <- article
	}
	return article, nil
}

func (repo *ArticleRepository) GetSingle(id string) (*model.Article, error) {
	item, err := repo.db.Value(id)
	if item == nil || err == cache2go.ErrKeyNotFound {
		return nil, errors.New(fmt.Sprintf("no article found with id \"%s\"", id))
	}
	return item.Data().(*model.Article), nil
}

func (repo *ArticleRepository) GetAll() ([]*model.Article, error) {
	var articles []*model.Article
	repo.db.Foreach(func(id interface{}, item *cache2go.CacheItem) {
		articles = append(articles, item.Data().(*model.Article))
	})
	return articles, nil
}

func (repo *ArticleRepository) Delete(id string) (bool, error) {
	_, err := repo.db.Delete(id)
	if err == cache2go.ErrKeyNotFound {
		return false, errors.New(fmt.Sprintf("no article found with id \"%s\"", id))
	}
	return true, nil
}

func (repo *ArticleRepository) GetCreationStream() chan *model.Article {
	return repo.creationStream
}

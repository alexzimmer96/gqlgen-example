package service

import (
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/alexzimmer96/gqlgen-example/repository"
)

type IArticleService interface {
	ListArticles() ([]*model.Article, error)
	GetArticle() (*model.Article, error)
	SaveArticle(*model.Article) (*model.Article, error)
	DeleteArticle(string) (bool, error)
	CreateArticleFromRequest(*model.UpdateArticle) (*model.Article, error)
}

type ArticleService struct {
	repo repository.IArticleRepository
}

func NewArticleService(repo repository.IArticleRepository) *ArticleService {
	return &ArticleService{repo}
}

func (svc *ArticleService) ListArticles() ([]*model.Article, error) {
	return svc.repo.GetAll()
}

func (svc *ArticleService) GetArticle(id string) (*model.Article, error) {
	return svc.repo.GetSingle(id)
}

func (svc *ArticleService) SaveArticle(article *model.Article) (*model.Article, error) {
	return svc.repo.Save(article)
}

func (svc *ArticleService) DeleteArticle(id string) (bool, error) {
	return svc.repo.Delete(id)
}

func (svc *ArticleService) CreateArticleFromRequest(creationRequest *model.UpdateArticle) (*model.Article, error) {
	parsedArticle, err := creationRequest.TransformToArticle()
	if err != nil {
		return nil, err
	}
	return svc.SaveArticle(parsedArticle)
}

func (svc *ArticleService) GetCreationStream() chan *model.Article {
	return svc.repo.GetCreationStream()
}

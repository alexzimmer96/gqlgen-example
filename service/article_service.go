package service

import (
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/alexzimmer96/gqlgen-example/repository"
)

type IArticleService interface {
	ListArticles() ([]*model.Article, error)
	GetArticle(string) (*model.Article, error)
	SaveArticle(*model.Article) (*model.Article, error)
	DeleteArticle(string) (bool, error)
	CreateArticleFromRequest(*model.CreateArticle) (*model.Article, error)
	ApplyArticleChanges(string, *model.UpdateArticle) (*model.Article, error)
	SubscribeArticleCreation() (chan *model.Article, chan bool)
}

type ArticleService struct {
	repo                     repository.IArticleRepository
	articleCreationObservers []*ArticleServiceObserver
}

type ArticleServiceObserver struct {
	CreationStream chan *model.Article
}

func NewArticleService(repo repository.IArticleRepository) *ArticleService {
	service := &ArticleService{repo: repo}
	go service.articleCreationStreamMultiplexer()
	return service
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

func (svc *ArticleService) CreateArticleFromRequest(creationRequest *model.CreateArticle) (*model.Article, error) {
	parsedArticle, err := creationRequest.TransformToArticle()
	if err != nil {
		return nil, err
	}
	return svc.SaveArticle(parsedArticle)
}

func (svc *ArticleService) ApplyArticleChanges(id string, changes *model.UpdateArticle) (*model.Article, error) {
	article, err := svc.repo.GetSingle(id)
	if err != nil {
		return nil, err
	}
	article, err = changes.MergeChanges(article)
	if err != nil {
		return nil, err
	}
	return svc.SaveArticle(article)
}

func (svc *ArticleService) SubscribeArticleCreation() *ArticleServiceObserver {
	deliveryChannel := make(chan *model.Article)
	observer := &ArticleServiceObserver{deliveryChannel}
	svc.articleCreationObservers = append(svc.articleCreationObservers, observer)
	return observer
}

func (svc *ArticleService) UnsubscribeArticleCreation(observer *ArticleServiceObserver) {
	j := 0
	for _, entry := range svc.articleCreationObservers {
		if entry == observer {
			svc.articleCreationObservers[j] = entry
		}
	}
	svc.articleCreationObservers = svc.articleCreationObservers[:j]
}

func (svc *ArticleService) articleCreationStreamMultiplexer() {
	incoming := svc.repo.GetCreationStream()
	for {
		select {
		case article := <-incoming:
			for _, entry := range svc.articleCreationObservers {
				entry.CreationStream <- article
			}
		}
	}
}

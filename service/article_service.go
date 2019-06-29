package service

import (
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/alexzimmer96/gqlgen-example/repository"
	"sync"
)

type IArticleService interface {
	ListArticles() ([]*model.Article, error)
	GetArticle(string) (*model.Article, error)
	SaveArticle(*model.Article) (*model.Article, error)
	DeleteArticle(string) (bool, error)
	CreateArticleFromRequest(*model.CreateArticle) (*model.Article, error)
	ApplyArticleChanges(string, *model.UpdateArticle) (*model.Article, error)
	SubscribeArticleCreation() (chan model.Article, chan bool)
}

type ArticleService struct {
	repo                     repository.IArticleRepository
	articleCreationObservers []*ArticleServiceObserver
	mutex                    sync.Mutex
}

type ArticleServiceObserver struct {
	CreationStream chan *model.Article
}

func NewArticleService(repo repository.IArticleRepository) *ArticleService {
	service := &ArticleService{repo: repo}
	go service.articleCreationStreamMultiplexer()
	return service
}

// Returns a List of all Articles found in the Database.
func (svc *ArticleService) ListArticles() ([]*model.Article, error) {
	return svc.repo.GetAll()
}

// Returns an Article from the Database by its ID.
func (svc *ArticleService) GetArticle(id string) (*model.Article, error) {
	return svc.repo.GetSingle(id)
}

// Saves an Article to the Database. If the Article does not already exists, it will be created.
func (svc *ArticleService) SaveArticle(article *model.Article) (*model.Article, error) {
	return svc.repo.Save(article)
}

// Delete an Article by its ID. Returning (true, nil) when the deletion was successful.
func (svc *ArticleService) DeleteArticle(id string) (bool, error) {
	return svc.repo.Delete(id)
}

// Creates a new Article from an CreateArticle-Request
func (svc *ArticleService) CreateArticleFromRequest(creationRequest *model.CreateArticle) (*model.Article, error) {
	parsedArticle, err := creationRequest.TransformToArticle()
	if err != nil {
		return nil, err
	}
	return svc.SaveArticle(parsedArticle)
}

// Applies a UpdateArticle-Request to an Article by its IDs and stores the updated Article back to the database.
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

// Adds an Observer to the ArticleCreationStream. The returned object is holding a personal channel
func (svc *ArticleService) SubscribeArticleCreation() *ArticleServiceObserver {
	svc.mutex.Lock()
	deliveryChannel := make(chan *model.Article)
	observer := &ArticleServiceObserver{deliveryChannel}
	svc.articleCreationObservers = append(svc.articleCreationObservers, observer)
	svc.mutex.Unlock()
	return observer
}

// Remove an Observer from the ArticleCreationStream
func (svc *ArticleService) UnsubscribeArticleCreation(observer *ArticleServiceObserver) {
	svc.mutex.Lock()
	close(observer.CreationStream)
	j := 0
	for _, entry := range svc.articleCreationObservers {
		if entry == observer {
			svc.articleCreationObservers[j] = entry
		}
	}
	svc.articleCreationObservers = svc.articleCreationObservers[:j]
	svc.mutex.Unlock()
}

// This multiplexer routes incoming Articles from the the articleRepository to every active subscriber.
func (svc *ArticleService) articleCreationStreamMultiplexer() {
	incoming := svc.repo.GetCreationStream()
	for {
		article := <-incoming
		svc.mutex.Lock()
		for _, entry := range svc.articleCreationObservers {
			entry.CreationStream <- article
		}
		svc.mutex.Unlock()
	}
}

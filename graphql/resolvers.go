package graphql

import (
	"context"
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/alexzimmer96/gqlgen-example/service"
	"github.com/prometheus/client_golang/prometheus"
)

type Resolver struct {
	articleService *service.ArticleService
}

func NewResolver(articleService *service.ArticleService) *Resolver {
	return &Resolver{
		articleService: articleService,
	}
}

//======================================================================================================================

var (
	resolverRequest        = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "graphql_resolver_requests"}, []string{"action_type", "action"})
	subscriptionDeliveries = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "graphql_subscription_deliveries"}, []string{"action"})
)

func init() {
	prometheus.MustRegister(resolverRequest, subscriptionDeliveries)
}

//======================================================================================================================

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	resolverRequest.WithLabelValues("query", "Articles").Add(1)
	return r.articleService.ListArticles()
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	resolverRequest.WithLabelValues("query", "Article").Add(1)
	return r.articleService.GetArticle(id)
}

//======================================================================================================================

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateArticle(ctx context.Context, article model.CreateArticle) (*model.Article, error) {
	resolverRequest.WithLabelValues("mutation", "CreateArticle").Add(1)
	return r.articleService.CreateArticleFromRequest(&article)
}

func (r *mutationResolver) UpdateArticle(ctx context.Context, id string, update model.UpdateArticle) (*model.Article, error) {
	resolverRequest.WithLabelValues("mutation", "UpdateArticle").Add(1)
	return r.articleService.ApplyArticleChanges(id, &update)
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (bool, error) {
	resolverRequest.WithLabelValues("mutation", "DeleteArticle").Add(1)
	return r.articleService.DeleteArticle(id)
}

//======================================================================================================================

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) ArticleCreated(ctx context.Context) (<-chan *model.Article, error) {
	subscription := r.articleService.SubscribeArticleCreation()
	go func() {
		<-ctx.Done()
		r.articleService.UnsubscribeArticleCreation(subscription)
	}()
	return subscription.CreationStream, nil
}

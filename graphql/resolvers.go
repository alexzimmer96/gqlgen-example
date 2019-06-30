package graphql

import (
	"context"
	"github.com/alexzimmer96/gqlgen-example/model"
	"github.com/alexzimmer96/gqlgen-example/service"
)

type Resolver struct {
	articleService service.IArticleService
}

func NewResolver(articleService service.IArticleService) *Resolver {
	return &Resolver{
		articleService: articleService,
	}
}

//======================================================================================================================

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	return r.articleService.ListArticles()
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	return r.articleService.GetArticle(id)
}

//======================================================================================================================

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateArticle(ctx context.Context, article model.CreateArticle) (*model.Article, error) {
	return r.articleService.CreateArticleFromRequest(&article)
}

func (r *mutationResolver) UpdateArticle(ctx context.Context, id string, update model.UpdateArticle) (*model.Article, error) {
	return r.articleService.ApplyArticleChanges(id, &update)
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (bool, error) {
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

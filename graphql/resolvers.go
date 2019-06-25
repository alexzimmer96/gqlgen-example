package graphql

import (
	"context"
	"github.com/alexzimmer96/gqlgen-example/model"
)

type Resolver struct {}

func NewResolver() *Resolver {
	return &Resolver{}
}

//======================================================================================================================

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	panic("implement me")
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	panic("implement me")
}

//======================================================================================================================

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateArticle(ctx context.Context, article model.UpdateArticle) (*model.Article, error) {
	panic("implement me")
}

func (r *mutationResolver) UpdateArticle(ctx context.Context, id string, update model.UpdateArticle) (*model.Article, error) {
	panic("implement me")
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (bool, error) {
	panic("implement me")
}

//======================================================================================================================

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) ArticleCreated(ctx context.Context) (<-chan *model.Article, error) {
	panic("implement me")
}


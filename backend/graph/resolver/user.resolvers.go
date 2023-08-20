package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"backend/graph"
	"backend/graph/model"
	"backend/module"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.UserPayload, error) {
	id := uuid.New().String()
	if input.ScreenName == nil {
		rG := module.GenerateRandom{}
		rG.Init()
		screenName := rG.GetAlphanumberic(rG.GetRandom(15, 3))
		input.ScreenName = &screenName
	}

	var posts []*model.Post

	user := model.User{
		ID:         id,
		UserName:   input.UserName,
		ScreenName: *input.ScreenName,
		Posts:      posts,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return &model.UserPayload{
			Error: &model.Error{
				Message: err.Error(),
			},
		}, err
	}
	return &model.UserPayload{
		User: &user,
	}, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.UserPayload, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, input model.ModelInputID) (*model.DeleteUserPayload, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - deleteUser"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, input model.ModelInputID) (*model.GetUserPayload, error) {
	var user model.User
	err := r.DB.First(&user, input.ID).Error

	return &model.GetUserPayload{
		User: &user,
	}, err
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, input model.ConnectionInput) (*model.GetUsersPayload, error) {
	var users []*model.User
	if err := r.DB.Find(&users).Error; err != nil {
		return &model.GetUsersPayload{
			Error: &model.Error{
				Message: err.Error(),
			},
		}, err
	}

	return &model.GetUsersPayload{
		Users: users,
	}, nil
}

// Posts is the resolver for the posts field.
func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
	var posts []*model.Post
	err := r.DB.Find(&posts, model.Post{AuthorID: obj.ID}).Error
	if module.IsErrorExcludeNoneRecord(err) {
		return nil, err
	}
	return posts, nil
}

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

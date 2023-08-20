package loader

import (
	"backend/graph/model"
	"backend/module"
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

type PostLoader struct {
	DB *gorm.DB
}

func (u *PostLoader) BatchGetPosts(ctx context.Context, authorIDs []string) []*dataloader.Result[[]*model.Post] {
	postsTemp := []*model.Post{}
	if err := u.DB.Find(&postsTemp, "author_id IN ?", authorIDs).Error; module.IsErrorExcludeNoneRecord(err) {
		err := fmt.Errorf("error, %w", err)
		log.Printf("%v\n", err)
		panic(err)
	}

	postsByAuthorId := map[string][]*model.Post{}
	for _, post := range postsTemp {
		postsByAuthorId[post.AuthorID] = append(postsByAuthorId[post.AuthorID], post)
	}

	result := make([]*dataloader.Result[[]*model.Post], len(authorIDs))
	for index, ID := range authorIDs {
		post := postsByAuthorId[ID]
		result[index] = &dataloader.Result[[]*model.Post]{Data: post, Error: nil}
	}
	return result
}

func LoadPosts(ctx context.Context, authorId string) ([]*model.Post, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.PostLoader.Load(ctx, authorId)
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result, nil
}

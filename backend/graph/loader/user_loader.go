package loader

import (
	"backend/graph/model"
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

type UserLoader struct {
	DB *gorm.DB
}

func (u *UserLoader) BatchGetUsers(ctx context.Context, userIDs []string) []*dataloader.Result[*model.User] {
	usersTemp := []*model.User{}
	if err := u.DB.Find(&usersTemp, "id IN ?", userIDs).Error; err != nil {
		err := fmt.Errorf("fail get users, %w", err)
		log.Printf("%v\n", err)
		panic(err)
	}

	usersByUserId := map[string]*model.User{}
	for _, user := range usersTemp {
		usersByUserId[user.ID] = user
	}

	users := make([]*model.User, len(userIDs))
	for i, id := range userIDs {
		users[i] = usersByUserId[id]
	}

	result := make([]*dataloader.Result[*model.User], len(userIDs))
	for index := range userIDs {
		user := users[index]
		result[index] = &dataloader.Result[*model.User]{Data: user, Error: nil}
	}
	return result
}

func LoadUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.UserLoader.Load(ctx, userID)
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result, nil
}

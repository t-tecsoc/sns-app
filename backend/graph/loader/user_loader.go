package loader

import (
	"backend/graph/model"
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader"
	"gorm.io/gorm"
)

type UserLoader struct {
	DB *gorm.DB
}

func (u *UserLoader) BatchGetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	userIDs := make([]string, len(keys))
	for i, key := range keys {
		userIDs[i] = key.String()
	}
	usersTemp := []*model.User{}
	if err := u.DB.Find(&usersTemp, "id IN ?", userIDs).Error; err != nil {
		err := fmt.Errorf("fail get users, %w", err)
		log.Printf("%v\n", err)
		return nil
	}

	usersByUserId := map[string]*model.User{}
	for _, user := range usersTemp {
		usersByUserId[user.ID] = user
	}

	users := make([]*model.User, len(userIDs))
	for i, id := range userIDs {
		users[i] = usersByUserId[id]
	}

	output := make([]*dataloader.Result, len(userIDs))
	for index := range userIDs {
		user := users[index]
		output[index] = &dataloader.Result{Data: user, Error: nil}
	}
	return output
}

// dataloader.Loadをwrapして型づけした実装
func LoadUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

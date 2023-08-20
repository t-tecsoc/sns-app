package loader

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader"
	"gorm.io/gorm"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders 各DataLoaderを取りまとめるstruct
type Loaders struct {
	UserLoader *dataloader.Loader
}

func NewLoaders(db *gorm.DB) *Loaders {
	// define the data loader
	userLoader := &UserLoader{
		DB: db,
	}
	loaders := &Loaders{
		UserLoader: dataloader.NewBatchedLoader(
			userLoader.BatchGetUsers,
			dataloader.WithClearCacheOnBatch(),
		),
	}
	return loaders
}

// Middleware LoadersをcontextにインジェクトするHTTPミドルウェア
func Middleware(loaders *Loaders, next gin.HandlerFunc) gin.HandlerFunc {
	loaders.UserLoader.ClearAll()
	// return a middleware that injects the loader to the request context
	return func(c *gin.Context) {
		cCtx := c.Request.Context()
		ctx := context.WithValue(cCtx, loadersKey, loaders)
		c.Request = c.Request.WithContext(ctx)
		next(c)
	}
}

// GetLoaders ContextからLoadersを取得する
func GetLoaders(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

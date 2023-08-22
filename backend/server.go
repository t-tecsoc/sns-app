package main

import (
	"backend/db"
	"backend/graph"
	"backend/graph/loader"
	"backend/graph/model"
	"backend/graph/resolver"
	"backend/graph/validation"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Graphql handler
func graphqlHandler(config graph.Config) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	h := handler.NewDefaultServer(
		graph.NewExecutableSchema(config),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func loadEnv() error {
	err := godotenv.Load("../.env")
	return err
}

func main() {

	err := loadEnv()

	if err != nil {
		fmt.Printf(".envファイルを読み込み出来ませんでした\n %v", err)
		return
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST_NAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DB_CONNECTION_PORT"),
	)

	DB := db.ConnectGORM(dsn)
	// loaderの初期化
	ldrs := loader.NewLoaders(DB)
	DB.Logger = DB.Logger.LogMode(logger.Info)

	DB.AutoMigrate(&model.User{}, &model.Post{})

	config := graph.Config{
		Resolvers: &resolver.Resolver{
			DB: DB,
		},
	}
	config.Directives.Validation = func(ctx context.Context, obj interface{}, next graphql.Resolver, format string) (res interface{}, err error) {
		errors, err := validation.ValidateModel(ctx)
		if err != nil {
			return nil, err
		}

		if len(errors) > 0 {
			log := ""
			for _, e := range errors {
				// 改行を入れてlogに追加する
				log += e + "\n"
			}
			return nil, fmt.Errorf(log)
		}

		return next(ctx)
	}

	// Setting up Gin
	router := gin.Default()
	router.Use(GinContextToContextMiddleware())
	srv := graphqlHandler(config)
	router.POST("/query", loader.Middleware(ldrs, srv))

	if gin.Mode() != gin.ReleaseMode {
		router.GET("/", playgroundHandler())
	}
	var addr = fmt.Sprintf("%s:%s", os.Getenv("BACKEND_HOSTNAME"), os.Getenv("BACKEND_CONTAINER_PORT"))
	router.Run(addr)
}

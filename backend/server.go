package main

import (
	"backend/db"
	"backend/graph"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Graphql handler
func graphqlHandler(db *gorm.DB) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					// resolver.goで宣言した構造体にデータベースの値を受け渡し
					DB: db,
				},
			},
		),
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

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func loadEnv() error {
	err := godotenv.Load("../.env")
	return err
}

func main() {
	err := loadEnv()
	if err != nil {
		fmt.Printf(".envファイルを読み込み出来ませんでした: %v", err)
		return
	}
	database_url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("HOST_NAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"),
	)

	db := db.ConnectGORM(database_url)

	// resolver内でデータベースを扱えるように設定

	// Setting up Gin
	r := gin.Default()
	r.Use(GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler(db))

	if gin.Mode() != gin.ReleaseMode {
		r.GET("/", playgroundHandler())
	}

	r.Run()
}

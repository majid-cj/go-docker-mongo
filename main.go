package main

import (
	"fmt"
	"log"
	"os"

	"github.com/majid-cj/go-docker-mongo/infrastructure/auth"
	"github.com/majid-cj/go-docker-mongo/infrastructure/persistence"
	"github.com/majid-cj/go-docker-mongo/router"
	"github.com/majid-cj/go-docker-mongo/util/middleware"

	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}
	os.Mkdir(os.Getenv("UPLOADS"), os.FileMode(0766))
}

func main() {
	repositoryservice, err := persistence.NewRepository()
	authservice := auth.NewDBAuth()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer repositoryservice.Client.Disconnect(repositoryservice.Ctx)
	defer authservice.DB.Close()

	CORS := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(logger.New())
	app.Use(recover.New())
	app.UseRouter(CORS)
	app.UseGlobal(middleware.RateLimit)

	app.HandleDir("/uploads", iris.Dir("./uploads"))
	app.I18n.Load("./locales/*/*.yaml")
	app.I18n.SetDefault("en")

	router.API(app, repositoryservice, authservice)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")), iris.WithOptimizations)
}

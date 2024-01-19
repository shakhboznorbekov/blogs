package main

import (
	"fmt"
	auth2 "github.com/blogs/internal/auth"
	"github.com/blogs/internal/controller/http/v1/auth"
	blog_controller "github.com/blogs/internal/controller/http/v1/blog"
	new_controller "github.com/blogs/internal/controller/http/v1/new"
	user_controller "github.com/blogs/internal/controller/http/v1/user"
	"github.com/blogs/internal/pkg/config"
	"github.com/blogs/internal/pkg/repository/postgres"
	"github.com/blogs/internal/pkg/script"
	blog_repo "github.com/blogs/internal/repository/postgres/blog"
	new_repo "github.com/blogs/internal/repository/postgres/new"
	user_repo "github.com/blogs/internal/repository/postgres/user"
	"log"

	"github.com/blogs/internal/router"
)

func main() {
	// config
	cfg := config.GetConf()

	// databases
	postgresDB := postgres.New(cfg.DBUsername, cfg.DBPassword, cfg.DBPort, cfg.DBName, config.GetConf().DefaultLang, config.GetConf().BaseUrl)

	//migration
	script.MigrateUP(postgresDB)

	// authenticator
	authenticator := auth2.New(postgresDB)

	//repository
	userRepo := user_repo.NewRepository(postgresDB)
	newRepo := new_repo.NewRepository(postgresDB)
	blogRepo := blog_repo.NewRepository(postgresDB)

	//controller
	userController := user_controller.NewController(userRepo, authenticator)
	authController := auth.NewController(userRepo, authenticator)
	newController := new_controller.NewController(newRepo)
	blogController := blog_controller.NewController(blogRepo)

	// router
	r := router.New(authenticator, userController, authController, blogController, newController)
	log.Fatalln(r.Init(fmt.Sprintf(":%s", cfg.Port)))

}

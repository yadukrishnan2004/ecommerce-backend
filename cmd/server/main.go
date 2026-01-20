package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/notifications"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repositery"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/infrastructure"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/router"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/service"
)

func main() {

	//initilizing the fiber router
	app:=fiber.New()

	//loading the env , setting the Port for runnint the server, setting the DSN 
	cfg:=config.Load()

	// connecting the data base pass an dsn (data source name in the form of string)
	DB:=infrastructure.ConnectPostgres(cfg.DSN)

	//Auto migrate repositery tables 

	DB.AutoMigrate(
		&domain.User{},   //user table
	)

	// setting up the handler layer


	nofier:=notifications.NewemailNodifier(
									cfg.SMTP_HOST,
									cfg.SMTP_PORT,
									cfg.SMTP_EMAIL,
									cfg.SMTP_PASS,
								)

	//Reopsiterys
	userRepo:=repositery.NewUserRepo(DB)

	//services
	userSVC:=service.NewUserService(userRepo,nofier)

	//handlers
	UserHandler:=handler.NewUserHandler(userSVC)
	
	//setting up the router 

	router.SetUpRouther(app,UserHandler)

	app.Listen(":"+cfg.App_Port)

}
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	auth "github.com/yadukrishnan2004/ecommerce-backend/internal/Auth"
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

	
	//loading the env , setting the Port for runnint the server, setting the DSN 
	cfg:=config.Load()
	
	// connecting the data base pass an dsn (data source name in the form of string)
	DB:=infrastructure.ConnectPostgres(cfg.DSN)
	sqlDb,err:=DB.DB()
	if err != nil {
		log.Fatal("fale to underlaying db connection")
	}
	
	//initilizing the fiber router
	app:=fiber.New(fiber.Config{
			DisableStartupMessage: true, 
		})
		
		
		//Auto migrate repositery tables 
		
		DB.AutoMigrate(
			&domain.User{},   //user table
		)
	
	
	// setting up the handler layer


	nofier:=notifications.NewemailNodifier( cfg.SMTP_HOST,
											cfg.SMTP_PORT,
											cfg.SMTP_EMAIL,
											cfg.SMTP_PASS)

	//Reopsiterys
	userRepo:=repositery.NewUserRepo(DB)


	//jwt service 
	jwt:=auth.NewJwtService(cfg.JWT)

	//services
	userSVC:=service.NewUserService(userRepo,nofier,*jwt)
	AdminSVC:=service.NewAdminService(userRepo)

	//handlers
	UserHandler:=handler.NewUserHandler(userSVC)
	AdminHandler:=handler.NewAdminHandler(AdminSVC)
	
	//setting up the router 

	router.SetUpRouther(app,UserHandler,AdminHandler)

//runing the server in an separate goroutine
	go func(){
		fmt.Printf("Server is running on :%s\n",cfg.App_Port)
		if err:=app.Listen(":"+cfg.App_Port);err != nil {
			log.Panic(err)
		}
	}()

	c:=make(chan os.Signal,1)

	signal.Notify(c,os.Interrupt,syscall.SIGTERM)
	<-c
	fmt.Println("The Sever is Starting shutting down in 5 seconds....")

		ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		
		if err:=app.ShutdownWithContext(ctx);err !=nil{
			fmt.Printf("Server force to Shutdown:%v\n",err)
		}

		fmt.Println("Closing data base connection...")
		if err:=sqlDb.Close();err != nil {
			fmt.Printf("Faile to Close the db:%v\n",err)

		}
		fmt.Println("Sever shutdown successfully")
}
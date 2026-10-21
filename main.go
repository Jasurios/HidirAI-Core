package main

import (
	"HidirAi-Core/database"
	"HidirAi-Core/handlers"
	"HidirAi-Core/middleware"
	"log"
	"os"

	logger "github.com/Jasurios/SimpleLogger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main()  {
	err := logger.Start("server.log")
	if err != nil{
		log.Fatalln(err)
	}
	log.Println("Logger started")

	err = godotenv.Load("./config.env")
	if err != nil{
		log.Fatalln(err)
	}
	log.Println("Vars loaded")

	db , err := database.Open(os.Getenv("DATABASE"))
	if err != nil{
		log.Fatalln(err)
	}
	err = database.Init(db)
	if err != nil{
		log.Fatalln(err)
	}
	db.Close()

	server := gin.Default()

	server.POST("/newkey",middleware.SystemAuth(os.Getenv("ADMIN_PASSWORD")), handlers.Newkey)
	server.POST("/ban/:user",middleware.SystemAuth(os.Getenv("ADMIN_PASSWORD")), handlers.Ban)
	server.POST("/del/:user",middleware.SystemAuth(os.Getenv("ADMIN_PASSWORD")), handlers.DelUser)

	server.POST("/chat",middleware.UserAuth(),handlers.Chat)
	server.POST("/service",middleware.ServiceAuth(),handlers.ChatService)

	server.POST("/newusermcp",middleware.UserAuth(),handlers.NewUserMCP)
	server.POST("/newservicemcp",middleware.ServiceAuth(),handlers.NewServiceMCP)

	server.Run(":2010")
}
package handlers

import (
	"HidirAi-Core/database"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Ban(c *gin.Context)  {
	user := c.Param("user")
	log.Println(user)
	db, err := database.Open(os.Getenv("DATABASE"))
	if err != nil{
		log.Println(err)
		return
	}

	banned, err := database.Ban(db, user)
	if err != nil{
		log.Println(err)
		return
	}
	defer db.Close()

	if banned{
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.AbortWithStatus(http.StatusAlreadyReported)
}
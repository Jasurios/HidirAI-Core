package handlers

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func DelUser(c *gin.Context)  {
	user := c.Param("user")

	err := os.RemoveAll("./users/"+user)
	if err != nil{
		log.Println(err)
	}
}
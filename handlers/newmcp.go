package handlers

import (
	"HidirAi-Core/structs"
	"HidirAi-Core/wwf"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUserMCP(c *gin.Context)  {
	user := c.GetString("user")
	var server structs.MCPServer

	err := c.ShouldBindJSON(&server)
	if err != nil || server.Url == ""{
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = wwf.Newmcp(user, server.Url, server.Key)
	if err != nil{
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusOK)
}

func NewServiceMCP(c *gin.Context)  {
	var server structs.MCPServerService

	err := c.ShouldBindJSON(&server)
	if err != nil || server.User == "" || server.Url == ""{
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = wwf.Newmcp(server.User, server.Url, server.Key)
	if err != nil{
		log.Println(err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

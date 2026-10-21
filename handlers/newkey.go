package handlers

import (
	"HidirAi-Core/database"
	"HidirAi-Core/structs"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func genapi() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "hid_" + hex.EncodeToString(b), nil
}

func Newkey(c *gin.Context){
	var key structs.Newkey

	err := c.ShouldBindJSON(&key)
	if err != nil || key.Type != "Service" && key.Type != "User" || key.User == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return 
	}

	api, err := genapi()

	db, err := database.Open(os.Getenv("DATABASE"))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	
	err = database.SaveKey(db, key.Type, key.User, api)
	if err != nil{
		log.Println(err)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"key":api})
}

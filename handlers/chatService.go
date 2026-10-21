package handlers

import (
	"HidirAi-Core/ai"
	"HidirAi-Core/database"
	"HidirAi-Core/structs"
	"HidirAi-Core/wwf"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ChatService(c *gin.Context) {
	var chat structs.ChatService

	err := c.ShouldBindJSON(&chat)
	if err != nil || chat.User == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	db, err := database.Open(os.Getenv("DATABASE"))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	keytype, err := database.KeyType(db, c.GetString("key"))
	if keytype != "Service" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	banned, err := database.CheckUser(db, chat.User)
	if banned {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Banned"})
		return
	}

	platform, _ := database.KeyOwner(db, c.GetString("key"))
	userpromt := ai.Getprompt(platform)+chat.Sysprompt

	err = wwf.CheckSpace(chat.User)
	if err != nil{
		log.Println(err)
	}
	
	var result ai.SkillResult

	if chat.Audio != nil{
		text, err := ai.STT(chat.Audio)
		if err != nil{
			log.Println(err)
		}
		result, err = ai.HidirAi(chat.User, text, userpromt)
		if err != nil{
			log.Println(err)
		}
	}else if chat.Image != ""{
		text, err := ai.ViewImage(chat.User, chat.Image, "photo")
		if err != nil{
			log.Println(err)
		}
		result, err = ai.HidirAi(chat.User, text, userpromt)
		if err != nil{
			log.Println(err)
		}
	}else if chat.Text != ""{
		result, err = ai.HidirAi(chat.User, chat.Text, userpromt)
		if err != nil{
			log.Println(err)
		}
	}

	switch result.Type {
	case ai.ResultText:
		resp := structs.ChatResponse{
			Text: result.Text,
		}
		c.JSON(http.StatusOK, resp)
	case ai.ResultPhoto:
		resp := structs.ChatResponse{
			Image: result.Photo,
		}
		c.JSON(http.StatusOK, resp)
	}
}

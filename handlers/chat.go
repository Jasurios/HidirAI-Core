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

func Chat(c *gin.Context) {
	var chat structs.ChatUser

	err := c.ShouldBindJSON(&chat)
	if err != nil || chat.Text == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	user := c.GetString("user")

	db, err := database.Open(os.Getenv("DATABASE"))
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	err = wwf.CheckSpace(user)
	if err != nil{
		log.Println(err)
	}
	var result ai.SkillResult

	if chat.Audio != nil {
		text, err := ai.STT(chat.Audio)
		if err != nil {
			log.Println(err)
		}
		result, err = ai.HidirAi(user, text, chat.Sysprompt)
		if err != nil {
			log.Println(err)
		}
	} else if chat.Image != "" {
		text, err := ai.ViewImage(user, chat.Image, "photo")
		if err != nil {
			log.Println(err)
		}
		result, err = ai.HidirAi(user, text, chat.Sysprompt)
		if err != nil {
			log.Println(err)
		}
	} else if chat.Text != "" {
		result, err = ai.HidirAi(user, chat.Text, chat.Sysprompt)
		if err != nil {
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

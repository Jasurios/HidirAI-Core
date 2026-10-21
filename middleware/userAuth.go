package middleware

import (
	"HidirAi-Core/database"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		
		
		db, err := database.Open(os.Getenv("DATABASE"))
		if err != nil{
			log.Println(err)
			return 
		}
		defer db.Close()

		user, _ := database.KeyOwner(db, key)
		c.Set("user", user)

		valid, err := database.CheckKey(db, key)
		if err != nil{
			log.Println(err)
			return 
		}
		if !valid{
			c.AbortWithStatus(http.StatusUnauthorized)
			return 
		}

		banned, err := database.CheckUser(db, user)
		if banned{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status":"Banned"})
			return
		}

		keytype, err := database.KeyType(db, key)
		if keytype != "User"{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
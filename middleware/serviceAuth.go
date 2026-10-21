package middleware

import (
	"HidirAi-Core/database"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ServiceAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		c.Set("key", key)

		db, err := database.Open(os.Getenv("DATABASE"))
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()

		valid, err := database.CheckKey(db, key)
		if err != nil {
			log.Println(err)
			return
		}
		if !valid {
			log.Println(valid)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

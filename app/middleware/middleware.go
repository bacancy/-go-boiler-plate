package middleware

import (
	"bacancy/go-boiler-plate/app/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("email", token.Email)
	}
}

func ValidateAdminToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"description": err.Error()})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Username or tokenString empty"})
			c.Abort()
			return
		}

		if token.Admin != true {
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Not admin detected"})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("email", token.Email)
	}
}

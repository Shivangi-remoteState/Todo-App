package middleware

import (
	"mytodoApp/database/dbHelper"
	"mytodoApp/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	get token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			c.Abort()
			return
		}
		//format
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&utils.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return utils.SECRET, nil
			},
			jwt.WithExpirationRequired(),
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		//if err := claims.Valid(); err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		//	c.Abort()
		//	return
		//}

		//	validate session
		_, err = dbHelper.GetUserIDBySessionID(claims.SessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("sessionID", claims.SessionID)

		c.Next()
	}
}

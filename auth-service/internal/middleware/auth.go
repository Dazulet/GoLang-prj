package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mangalib/auth-service/internal/utils"
)

const (
	keyUserID   = "userID"
	keyUserRole = "userRole"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(strings.TrimPrefix(header, "Bearer "), jwtSecret)
		if err != nil {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(keyUserID, claims.UserID)
		c.Set(keyUserRole, claims.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(keyUserRole)
		if role != "admin" {
			utils.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) uint {
	id, _ := c.Get(keyUserID)
	uid, _ := id.(uint)
	return uid
}

func CurrentUserRole(c *gin.Context) string {
	role, _ := c.Get(keyUserRole)
	r, _ := role.(string)
	return r
}

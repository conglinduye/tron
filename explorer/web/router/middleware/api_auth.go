package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/wlcy/tron/explorer/web/errno"
	"golang.org/x/crypto/bcrypt"
	"github.com/wlcy/tron/explorer/web/handler"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := ParseRequest(c); err != nil {
			fmt.Printf("AuthMiddleware err: %s\n", err)
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}


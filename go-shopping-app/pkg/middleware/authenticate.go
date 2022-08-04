package middleware

import (
	"fmt"
	"go-fruit-cart/pkg/apperrors"
	"go-fruit-cart/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			c.JSON(401, gin.H{"error": apperrors.ErrAuthHeaderMissing.Error()})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		fmt.Println("claims", claims)
		if err != nil {
			c.JSON(401, gin.H{"error": apperrors.ErrInvalidAuthToken.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims.Email)
		c.Next()
	}
}

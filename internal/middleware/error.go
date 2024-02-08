package middleware

import (
	"github.com/pkg/errors"
	"net/http"
	"xmicro/internal/httperror"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var e httperror.Http
			switch {
			case errors.As(err.Err, &e):
				c.AbortWithStatusJSON(e.Code, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}

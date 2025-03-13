package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func noRoute(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
}

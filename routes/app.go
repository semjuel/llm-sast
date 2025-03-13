package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/semjuel/llm-sast/handlers"
)

func app(r *gin.Engine) {
	r.POST("/api/app/upload/:model", handlers.Upload)
}

package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/semjuel/llm-sast/utils"
)

func Initialize(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowOrigins:    utils.GetAllowedOrigins(),
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Accept-Language", "Authorization"},
	}))

	r.OPTIONS("/*any", func(c *gin.Context) {
		return
	})

	noRoute(r)

	r.Static("/css", "./public/css")
	r.Static("/js", "./public/js")

	r.GET("/", func(c *gin.Context) {
		c.File("public/index.html")
	})

	// App file handling
	app(r)
}

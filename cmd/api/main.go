package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/semjuel/llm-sast/routes"
	"github.com/semjuel/llm-sast/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Subscribe to shut down signals; use cancel for premature shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = os.Stdout
	}

	r := gin.Default()
	proxies := utils.GetAllowedOrigins()
	if os.Getenv("APP_ENV") == "prod" && len(proxies) != 0 {
		//err = r.SetTrustedProxies(proxies)
	} else {
		err = r.SetTrustedProxies(nil)
	}

	if err != nil {
		panic(err)
	}

	r.TrustedPlatform = "X-Forwarded-For"

	routes.Initialize(r)

	go r.Run(":8099")

	// Waiting for shutdown signal.
	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

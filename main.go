package main

import (
	"GoCouchbase/api"
	"GoCouchbase/cluster"
	"GoCouchbase/config"
	"GoCouchbase/utils/logutil"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = config.GetConfig()
	_, _ = cluster.GetCluster()

	logger := logutil.GetLogger()
	logger.Println(logutil.Info, false, "Starting couchbase bridge server...")

	app := fiber.New()
	api.SetupRouter(app)

	logger.Println(logutil.Info, false, "Listening on port 8999...")
	err := app.Listen(":8999")
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Println(logutil.Info, false, "Shutting down couchbase bridge server...")
		os.Exit(0)
	}()
}

package main

import (
	"GoCouchbase/api"
	"GoCouchbase/cluster"
	"GoCouchbase/config"
	"GoCouchbase/utils/logutil"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	_ = config.GetConfig()
	_, _ = cluster.GetCluster()

	logger := logutil.GetLogger()
	logger.Println(logutil.Info, false, "Starting go couchbase server...")

	app := fiber.New()
	api.SetupRouter(app)

	err := app.Listen(":8999")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func createLogger() *log.Logger {
	logFile, err := os.OpenFile("/home/ubuntu/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}

func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "pong"})
}

type InventoryRequest struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Provider  string `json:"provider"`
	AccountID string `json:"accountId"`
	Service   string `json:"service"`
	Resource  string `json:"resource"`

	EndPoint string `json:"endPoint"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func GetRunEnv() string {
	runEnv := os.Getenv("RUN_ENV")
	return runEnv
}

func connectToCluster(env string, req *InventoryRequest) (*gocb.Cluster, error) {
	var connectionString, username, password string
	if env == "dev" || env == "prod" {
		connectionString = fmt.Sprintf("couchbase://%s", req.EndPoint)
		username = req.UserName
		password = req.Password
	} else {
		connectionString = "couchbase://localhost"
		username = "Administrator"
		password = "1234qwer!"
	}
	return gocb.Connect(connectionString, gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
}

func GetInventoryDetail(c *fiber.Ctx, logger *log.Logger) error {
	req := new(InventoryRequest)
	if err := c.BodyParser(req); err != nil {
		logger.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	runEnv := GetRunEnv()
	logger.Println("running env:", runEnv)
	logger.Println("request:", req)
	cluster, err := connectToCluster(runEnv, req)
	if err != nil {
		logger.Println("Error connecting to cluster:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	query := fmt.Sprintf("SELECT data.`%s` FROM %s WHERE provider='%s' AND account_id='%s' AND service='%s' AND resource='%s'", req.ID, req.TenantID, req.Provider, req.AccountID, req.Service, req.Resource)
	result, err := cluster.AnalyticsQuery(query, nil)
	if err != nil {
		logger.Println("Error executing query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for result.Next() {
		var row map[string]interface{}
		err := result.Row(&row)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Println("query result:", row)
		return c.JSON(row)
	}

	if err := result.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "can't find data with given parameters"})
}

func main() {
	logger := createLogger()

	app := fiber.New()
	app.Get("/ping", Ping)
	app.Post("/query", func(c *fiber.Ctx) error {
		return GetInventoryDetail(c, logger)
	})

	err := app.Listen(":8999")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

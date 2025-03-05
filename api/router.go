package api

import (
	"GoCouchbase/cluster"
	"GoCouchbase/inventory"
	"GoCouchbase/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	app.Get("/ping", Ping)
	app.Post("/query", func(c *fiber.Ctx) error {
		return getInventoryDetail(c)
	})
}

func Ping(c *fiber.Ctx) error {
	logger := utils.GetInfoLogger()
	logger.Println("ping")
	return c.JSON(fiber.Map{"message": "pong"})
}

func getInventoryDetail(c *fiber.Ctx) error {
	logger := utils.GetInfoLogger()
	errLogger := utils.GetErrorLogger()

	req := new(inventory.Request)
	if err := c.BodyParser(req); err != nil {
		errLogger.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	reqId := utils.GenerateUUID()
	req.RequestID = reqId
	logger.Printf("%+v\n", req)

	clusterInstance, err := cluster.GetCluster()
	if err != nil {
		errLogger.Println(reqId, "Error connecting to cluster:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	query := fmt.Sprintf("SELECT data.`%s` FROM %s WHERE provider='%s' AND account_id='%s' AND service='%s' AND resource='%s'", req.ID, req.TenantID, req.Provider, req.AccountID, req.Service, req.Resource)
	result, err := clusterInstance.AnalyticsQuery(query, nil)
	if err != nil {
		errLogger.Println(reqId, "Error executing query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for result.Next() {
		var row map[string]interface{}
		err = result.Row(&row)
		if err != nil {
			errLogger.Println(reqId, "Error parsing query result:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Println(reqId, "query result:", row)
		return c.JSON(row)
	}

	if err = result.Err(); err != nil {
		errLogger.Println(reqId, "Error executing query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	errLogger.Println(reqId, "No data found")
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "can't find data with given parameters"})
}

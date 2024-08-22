package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var NODE_NAME string
var NEXT_NODE string
var tracer trace.Tracer

func init() {
	nodeName, ok := os.LookupEnv("NODE_NAME")
	if !ok {
		nodeName = "application"
	}

	NODE_NAME = nodeName
	NEXT_NODE = os.Getenv("NEXT_NODE")
	tracer = otel.Tracer(NODE_NAME)

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	tp := initTracer()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := fiber.New()
	app.Use(otelfiber.Middleware())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"success": true,
		})
		return nil
	})

	app.Get("/api/v1/run", handleRunRequest)

	log.Fatal(app.Listen(":3000"))
}

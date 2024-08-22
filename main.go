package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var NODE_NAME string
var NEXT_NODE string
var tracer trace.Tracer
var client http.Client

func init() {
	nodeName, ok := os.LookupEnv("NODE_NAME")
	if !ok {
		nodeName = "application"
	}

	NODE_NAME = nodeName
	NEXT_NODE = os.Getenv("NEXT_NODE")
	tracer = otel.Tracer(NODE_NAME)

	client = http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func main() {
	tp := initTracer()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config := slogfiber.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}

	app := fiber.New()
	app.Use(otelfiber.Middleware(
		otelfiber.WithPort(3000),
		otelfiber.WithServerName(NODE_NAME),
		otelfiber.WithNext(func(ctx *fiber.Ctx) bool {
			return ctx.Path() == "/healthz"
		}),
	))

	app.Use(slogfiber.NewWithConfig(logger, config))

	app.Get("/healthz", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"success": true,
		})
		return nil
	})

	app.Get("/api/v1/run", handleRunRequest)

	log.Fatal(app.Listen(":3000"))
}

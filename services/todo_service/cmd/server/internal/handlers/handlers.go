package handlers

import (
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	goservice "github.com/onpointvn/libs/go-sdk"
	"github.com/onpointvn/libs/go-sdk/plugin/middleware"
)

func ping() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(&fiber.Map{
			"msg": "pong",
		})
	}
}

func Router(sc goservice.ServiceContext) *fiber.App {
	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})

	app.Use(logger.New(logger.Config{
		Format: `{"ip":${ip}, "timestamp":"${time}", "status":${status}, "latency":"${latency}", "method":"${method}", "path":"${path}"}` + "\n",
	}))
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.Recover(sc))
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	app.Get("/", ping())
	app.Get("/ping", ping())

	return app
}

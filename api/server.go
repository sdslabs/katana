package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	c "github.com/sdslabs/katana/api/controllers"
	cfg "github.com/sdslabs/katana/configs"
)

func Server() error {
	fiberConfig := fiber.Config{
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          30 * time.Second,
		DisableStartupMessage: true,
	}

	app := fiber.New(fiberConfig)
	app.Use(cors.New())

	corsConfig := cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Content-Length, Authorization, Cookie",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH",
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}

	app.Use(cors.New(corsConfig))

	api := app.Group("/api/v1")

	admin := api.Group("/admin")
	admin.Get("/:name", c.HelloAdmin)
	admin.Get("/cluster/:id", c.ClusterInfo)

	main := api.Group("/main")
	main.Get("/:name", c.HelloMain)

	fmt.Sprintf("%s:%d", cfg.APIConfig.Host, cfg.APIConfig.Port)
	return app.Listen(fmt.Sprintf("%s:%d", cfg.APIConfig.Host, cfg.APIConfig.Port))
}

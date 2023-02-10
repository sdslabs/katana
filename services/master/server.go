package master

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	cfg "github.com/sdslabs/katana/configs"
	c "github.com/sdslabs/katana/services/master/controllers"
)

func Server() error {
	fiberConfig := fiber.Config{
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          30 * time.Second,
		DisableStartupMessage: false,
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
	admin.Get("/echo/:name", c.HelloAdmin)
	admin.Get("/infraSet", c.InfraSet)
	admin.Get("/db", c.DB)
	admin.Post("/login", c.Login)
	admin.Post("/createTeams/:number", c.CreateTeams)
	admin.Post("/logs", c.Logs)
	admin.Get("/deploy", c.Deploy)
	admin.Get("/cluster/:id", c.ClusterInfo)
	fmt.Printf("Listening on %s:%d\n", cfg.APIConfig.Host, cfg.APIConfig.Port)
	return app.Listen(fmt.Sprintf("%s:%d", cfg.APIConfig.Host, cfg.APIConfig.Port))
}

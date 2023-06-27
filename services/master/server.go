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
		BodyLimit:             50 * 1024 * 1024,
		DisableStartupMessage: false,
	}

	app := fiber.New(fiberConfig)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

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
	admin.Get("/createTeams/:number", c.CreateTeams)
	admin.Post("/challengeUpdate", c.ChallengeUpdate)
	admin.Post("/logs", c.Logs)
	admin.Post("/deploy", c.Deploy)
	admin.Post("/deployCheckers", c.DeployCheckers)
	admin.Post("/updateChecker", c.UpdateChecker)
	admin.Get("/gitServer", c.GitServer)
	admin.Get("/cluster/:id", c.ClusterInfo)
	admin.Get("/deleteChallenge/:chall_name", c.DeleteChallenge)
	fmt.Printf("Listening on %s:%d\n", cfg.APIConfig.Host, cfg.APIConfig.Port)
	return app.Listen(fmt.Sprintf("%s:%d", cfg.APIConfig.Host, cfg.APIConfig.Port))
}

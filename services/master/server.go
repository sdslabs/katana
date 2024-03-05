package master

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	cfg "github.com/sdslabs/katana/configs"
	challengeDeployerService "github.com/sdslabs/katana/services/challengedeployerservice"
	infraSetService "github.com/sdslabs/katana/services/infrasetservice"
	c "github.com/sdslabs/katana/services/master/controllers"
)

func Server() error {
	fiberConfig := fiber.Config{
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          30 * time.Second,
		BodyLimit:             10 * 1024 * 1024,
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
	admin.Get("/infraSet", infraSetService.InfraSet)
	admin.Get("/db", infraSetService.DB)
	admin.Post("/login", infraSetService.Login)
	admin.Get("/createTeams", infraSetService.CreateTeams)
	admin.Post("/challengeUpdate", challengeDeployerService.ChallengeUpdate)
	admin.Get("/deploy", challengeDeployerService.Deploy)
	admin.Post("/deployChallenge", challengeDeployerService.DeployChallenge)
	admin.Get("/gitServer", infraSetService.GitServer)
	
	admin.Get("/cc",challengeDeployerService.Cc)
	admin.Get("/team",challengeDeployerService.Team)
	
	admin.Get("/deleteChallenge/:chall_name", challengeDeployerService.DeleteChallenge)
	log.Printf("Listening on %s:%d\n", cfg.APIConfig.Host, cfg.APIConfig.Port)
	return app.Listen(fmt.Sprintf("%s:%d", cfg.APIConfig.Host, cfg.APIConfig.Port))
}

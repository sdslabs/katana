package sensei

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/sdslabs/katana/api/services/sensei/controllers"
)

func ServiceInit() http.Handler {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	main := router.Group("/main")
	{
		main.GET("", c.HelloMain)
	}

	admin := router.Group("/admin")
	{
		admin.GET("", c.HelloAdmin)
	}

	return router
}

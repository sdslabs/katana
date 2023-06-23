package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DB(c *fiber.Ctx) error {
	client, err := utils.GetKubeClient("")
	if err != nil {
		log.Println(err)
	}
	mongo.Init()
	// Print the IP address of the service
	service, err := client.CoreV1().Services("katana").Get(context.TODO(), "mysql-nodeport-svc", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("MySQL Service IP:", service.Spec.ClusterIP)
	fmt.Println(service.Spec.Ports[0].Port) // NodePort (k8s service)
	//NodePort := strconv.Itoa(int(service.Spec.Ports[0].NodePort))
	mysql.Init()
	return c.SendString("Database setup completed\n")
}

func Login(c *fiber.Ctx) error {
	adminUser := new(types.AdminUser)

	if err := c.BodyParser(adminUser); err != nil {
		return err
	}

	admin, err := mongo.FetchSingleAdmin(adminUser.Username)
	if err != nil {
		return err
	}

	isAdmin := utils.CompareHashWithPassword(admin.Password, adminUser.Password)
	if isAdmin {
		jwtToken := jwt.New(jwt.SigningMethodHS256)

		token, err := jwtToken.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.HTTPOnly = true
		c.Cookie(cookie)

		return c.JSON(fiber.Map{
			"message": "success",
		})
	} else {
		return c.SendString("Incorrect password")
	}
}

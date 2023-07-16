package controllers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	hostsFilePath = "/etc/hosts"
)

func SetupIngress(c *fiber.Ctx) error {
	kubeClient, _ := utils.GetKubeClient()
	namespace := "kube-system"

	serviceName := "ingress-nginx-controller"
	if err := utils.WaitForLoadBalancerExternalIP(kubeClient, serviceName, namespace); err != nil {
		return err
	}

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	externalIP := service.Status.LoadBalancer.Ingress[0].IP

	hosts := []string{
		"mongo.katana.local",
		"gogs.katana.local",
		"mysql.katana.local",
		configs.KatanaConfig.Harbor.Hostname,
	}

	hostsEntry := ""
	for _, host := range hosts {
		hostsEntry += " " + host
	}

	file, err := os.Open(hostsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, hostsEntry) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				fields[0] = externalIP
				line = strings.Join(fields, " ")
				found = true
			}
		}
		lines = append(lines, line)
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s %s", externalIP, hostsEntry))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write to file
	file, err = os.OpenFile(hostsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	if w.Flush() != nil {
		return err
	} else {
		return c.JSON(fiber.Map{
			"message": "Successfully setup ingress",
		})
	}
}

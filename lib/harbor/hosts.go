package harbor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	config "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	hostsFilePath = "/etc/hosts"
)

var hostsEntry string = config.KatanaConfig.Harbor.Hostname

func checkHarborHostsEntryExists() bool {
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), hostsEntry) {
			return true
		}
	}

	return false
}

func addHarborHostsEntry() error {
	client, err := utils.GetKubeClient()
	if err != nil {
		return err
	}

	serviceName := "harbor"
	namespace := "katana"

	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	externalIP := service.Status.LoadBalancer.Ingress[0].IP

	print(externalIP)

	hostsEntry = fmt.Sprintf("%s %s", externalIP, hostsEntry)

	file, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := fmt.Fprintf(file, "\n%s\n", hostsEntry); err != nil {
		return err
	}

	return nil
}

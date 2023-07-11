package harbor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	config "github.com/sdslabs/katana/configs"
)

const (
	hostsFilePath = "/etc/hosts"
)

var hostsEntry string = "127.0.0.1 " + config.KatanaConfig.Harbor.Hostname

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

package controllers

import {
	"context"
	"fmt"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
}

 func logs(c *fiber.Ctx) error {

	//Loading kubeconfig file
	pathToCfg := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	fmt.Println(g.KatanaConfig.KubeNameSpace)

	// Load kubeclient
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	// Set up port forwarding
	forwarder, err := client.CoreV1().Pods("default").Forward(
		context.Background(),
		"es-cluster-0",
		kubernetes.NewPorts(9200, 9200),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer forwarder.Close()

	// Wait for the forwarding to start
	err = forwarder.Ready()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get("http://localhost:9200/_cluster/state?pretty")
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
	return c.SendString(string(body))
}
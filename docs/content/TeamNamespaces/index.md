---
title: "Team Namespaces"
---

As seen in the image below, each team has its own namespace. This is to ensure that each team has its own set of resources and does not interfere with other teams. This also allows for a more secure environment as each team can only access its own namespace.

![Image Not Found](/team-namespaces-architecture.png)

The following code in the ```createTeams.go``` controller is responsible for creating the namespaces for the teams and deploying the required resources into the namespaces.

```Golang
func CreateTeams(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }

	clusterConfig := g.ClusterConfig
	client, err := utils.GetKubeClient()
	if err != nil {
		log.Println(err)
	}
	noOfTeams, err := strconv.Atoi(c.Params("number"))

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < noOfTeams; i++ {
		log.Println("Creating Team: " + strconv.Itoa(i))
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "katana-team-ns-" + strconv.Itoa(i),
			},
		}
		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		manifest := &bytes.Buffer{}
		tmpl, err := template.ParseFiles(filepath.Join(clusterConfig.ManifestDir, "teams.yml"))
		if err != nil {
			return err
		}
		deploymentConfig := utils.DeploymentConfig()

		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}
		pathToCfg := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
		if err != nil {
			log.Fatal(err)
		}
		if err = deployment.ApplyManifest(config, client, manifest.Bytes(), "katana-team-ns-"+strconv.Itoa(i)); err != nil {
			return err
		}
	}
	return c.SendString("Successfully created teams")
}
```

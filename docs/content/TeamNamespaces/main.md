---
title: "Team Namespaces"
---

# Introduction

As seen in the image below, each team has its own namespace. This is to ensure that each team has its own set of resources and does not interfere with other teams. This also allows for a more secure environment as each team can only access its own namespace.

![Image Not Found](/team-namespaces-architecture.png)

# Setup

To deploy the database, you first need to set up the infrastructure with the help of the `/api/v2/admin/createTeams/:number` endpoint where `/:number` is the number of teams you want to initialise.

# Go Code For Team Namespaces

The following code in the `createTeams.go` controller is responsible for creating the namespaces for the teams and deploying the required resources into the namespaces.

This function can be found in the create teams controller, in the CreateTeams function.
```Golang
	for i := 0; i < noOfTeams; i++ {
		// Create a directory named katana-team-i in the teams directory
		if _, err := os.Stat("teams/katana-team-" + strconv.Itoa(i)); os.IsNotExist(err) {
			errDir := os.Mkdir("teams/katana-team-"+strconv.Itoa(i), 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}

		log.Println("Creating Team: " + strconv.Itoa(i))
		namespace := "katana-team-" + strconv.Itoa(i) + "-ns"
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
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
		if err = deployment.ApplyManifest(config, client, manifest.Bytes(), namespace); err != nil {
			return err
		}
	}
```

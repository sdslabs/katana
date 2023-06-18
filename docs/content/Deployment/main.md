---
title: "Deployment"
---
  

## Introduction

  

The deployment service is responsible for creating namespaces and deploying the ctf challenges for each team under their specific namespace along with the service attached with them.

  

{{<  toc  >}}

## Initialisation

- ### Create Teams
We need to make sure teams' namespaces are created and 
This creates the namespaces, the master pod for each pod. It also creates the user in Gogs database for each  team.

- ### Challenge Type
Currently we support web challenges for the ctf. An example of the web challenge can be found here.
https://github.com/dicegang/dicectf-2022-challenges/tree/master/web/notekeeper

- ### Challenge Zip
As an input to deployment service, we ask the zip file for the challenge.
TIP : make sure to create a zip without the path input flag for extra information.

## Sending Request
A post request to /deploy route under 'admin' Group with the type mulipart-form , key as challenge and the zip file as the load is required to be sent.
Here is an example of such a requst via postman.

![Image not found](/deploy-postreq.png)

You can also however send the request using frontend by setting up the saya frontend.

## Flow 

- ### Build Folders
Assuming the name of challenge is notekeeper for illustartion purposes. The file gets unzipped in ..../katana/chall/notekeeper along with a copy of the zip file in order to use for copying inside the pod. Read patching service for more info. 

- ### Build Image
Assuming that the docker file is in /notekeeper root base. Image is build outside in the docker registry or the volume mounted and pushed isnide minikube registry as of now. It is to be updated to use the docker client.

- ### Create Deployment
The deployment for each team with 1 replica of each challenge pod is created under each namespace for the teams using the k8's client. Check out deploy.go for the code file. 

- ### Create Service
Next, under each namespace, a Nodeport service is also attached to the deployments exposing the web challenge.Furthur exposing via minikube service command is also done.
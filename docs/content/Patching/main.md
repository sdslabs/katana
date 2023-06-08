---
title: "Patching Service"
---
  

## Introduction

  

The patching service of katana makes use of a locally run git service called Gogs running in the admin namespace. The decision to not use github was to decrease latency to pull stuff over internet. Following are the steps

  

{{<  toc  >}}

## Initialisation

- ### Infraset
During this time we estbalish MySQL, MongoDB and Gogs pods are created. 

- ### Database Setup
We establish connection with Mongo and MySQL. A mongoDB admin is established with team credentials. 

- ### GitServer
We hit the gogsIP/install which creates the gogs tables in the MySQL pod. If using a non-cloud based cluster (like minikube), establish a connection with LoadBalancer [```minikube tunnel```]. As of now you have to hit the Database setup one more time after GitServer to establish the admin user in Git 

- ### Create Teams
This creates the namespaces, the master pod for each pod. It also creates the user in Gogs database for each  team.

## Setting up a challenge

Whenever a challenge is setup, the broadcasting service is ivoked which creates a private repository for that challenge for each team along with applying a yaml file for that challenge in each team's namespace. Now the said broadcast service sends a zip file of the challenge to each and every pod, where it's unzipped and initialised dynamically w.r.t. each and every team repository. We pulll once to make sure the histories of both the local copy and the repository is in sync. 

## Patch Challenge
in /usr/bin/ we have a patch_challenge bash file which essentially runs a simple git add to git push command. It takes the commit message as it's argument so teams can have their own commit messages so they find it easier to point if they wish to backtrack to a previous patch. 
As soon as a push is made, a github webhook sends a post request upon which the updates are pulled, an image is created and pushed into the K8s registry. The challenge pod of that particular team is killed, and when it restarts, it pulls the latest image from the registry. 

## Inside/Out

To reduce latency, the following architectural decisions were taken. MySQL server and Gogs Service will be running within the K8s cluster. Pulling the changes and creation of image happens outside the cluster. 

## Image within

We "briefly" considered the concept of maybe creation of images within the pod and pushing it into the registry from within the pod itself. This would lead to the changes never having to leave the cluster, thus increasing speed by a lot. However, this was scrapped as to make an image you have to either: 

 - Create  a DIND (Docker in Docker), which would have secuirity impact. 
 - Using Kaniko (A Go based library which allows image creation within the pod). However pushing it into the registry recquired to provide sudo privilege to team pods, which essentially left the entire cluster vulnerable to attacks.
 
Thus the final decision was to run the patch service in the aforementioned manner to provide the team a seamlesss A&D  CTF exeperience and at the same time make it easier for the admin to host one

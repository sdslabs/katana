---
title: "SSH Service"
---

SSH Service enables team to remotely access their own mainpod.
Written by [Scar](https://github.com/Scar26), It is one of the fundamental components for the working of Katana. 

The setup the SSH service, run the following command (Considering you have katana running on your local system)

```Shell
curl localhost:3000/api/v1/admin/sshservice
```
This will run the SSH server on a `goroutine` (on a default port 2222) and generate a `teamcreds.txt` where you can fine the username and a randomly generated password for each team which can be used to access the pod. Once executed teams can access their masterpod in the following manner 
```Shell
ssh <team-username>@<host-ip> -p 2222
```

SSH Service comprises of the following functions, which are explained in detailed 

{{< toc >}}

## Init

This is basically the initialises `server.go`, getting the Configuration and Client set of the running Kubernetes Cluster and storing them in a global variable.

## Session Handler

This is the actual brains behind the entire thing. This essentially 'redirects' the ssh to a kubectl exec command, thus giving teams access to the Mainpod. 
 
## Password Handler
This is a basic function called whenever a team tries to ssh to make sure they have the correct password. If you want the access to be given regardless of the password provided, you can change 

`return utils.CompareHashWithPassword` 

to 

 `return true`

 `CompareHashWithPassword()` function essentially checks the password provided with the stored hash password (of course it's hashed, we can't be vulnerable ourselves while hosting an A&D CTF ;) ).  

## Server
This function has the basic purpose of returning the Server configuration for the SSH Service. This can then be executed with `ListenAndServe()` to run the server. The decision of not running the server directly upon invoking this server was done for ease of code, since this way config can be called and the server can start running later on, incase there is a requirement for the same  

## Create Teams

This function runs after `goroutine Server().ListenandServe()`. This gets the Labels of all katana-teams, generates a password for each of them and stores their Hash in the CTF Team datastructure, and prints the same in `teamcreds.txt` file.
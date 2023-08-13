---
title: "Tsuka"
---

Each team under it's own namespace has its own master pod. This master pod is Tsuka. This pod acts as a dedicated virtual machine for each team. It is a ubuntu based pod with ssh server running inside it.

All teams have access to their own master pod. They can ssh into it using the password provided to them. The password is stored in a file called teamcreds.txt which is generated when the team is created. NOTE: The teams have to be connected to Katana's VPN to access their master pod.

Tsuka contains source code of all the challenges. A team is expected to patch the challenge and push it to the git server. The patching service will then build the image and push it to the registry. The challenge pod will then pull the latest image and run it.

//UPDATE DIAGRAM OF A MASTER PODS WITH THINGS INSIDE IT.
![Image Not Found](/team-pods-architecture.png)

## Setup

During setup, a `setup-script.sh` is run which does the following:
- Does `apt-get update` and `apt-get upgrade` followed by `apt-get install git curl nano vim`.
- It then sets up the ssh server by installing openssh-server and setting up the password for the team. The password is the root password for the master pod as well.
- After this, it moves 2 binaries to `/usr/bin/` which are `patch-challenge` and `setup`. These binaries are used by the teams to patch their challenges and by the deployment service to setup challenges respectively.
- Lastly, it sets up git config for the team and runs a flask server as a daemon process.

## Working

- Tsuka contains a flask server which helps in setting up challenges when they are deployed. It runs as a daemon process and runs a `setup` script, which is present in `/usr/bin/`, when a challenge is deployed.

- There is also a `patch-challenge` binary present in `/usr/bin` that teams can use to patch their challenges. This binary is a wrapper around a bash script that runs a git add and git push command. The commit message is passed as an argument to the binary.

## General Flow

- Challenge gets copied to master pod
- Flask server unzips the challenge in the master pod
- Team makes changes to the challenge and uses `patch-challenge` binary to push the changes to the git server

For patching we looked at few options before finalising on the [current one](/Patching/). The other options were:

- [Challenge-Containers v1](/Tsuka/v1/)

- [Challenge-Containers v2](/Tsuka/v2)

# katana
An advanced yet simple attack/defence CTF infrastructure in Go

## Setup
- To start, you must have the following installed:
  - Go 1.18+
  - Minikube & kubectl
- Start a minikube cluster using the Docker driver with `minikube start --driver=docker`
- Enable the ingress addon with `minikube addons enable ingress`
- Apply the ingress-controller manifest with `kubectl apply -f manigests/expose-controller.yaml`
- Get the IP of the challenge deployer with `minikube service nginx-ingress-controller --url -n kube-system`
- Using the ip from the previous step, edit `/etc/hosts` and set local dns as `challengedeployer.katana.local` such as this `<ip>    challengedeployer.katana.local`. 

# katana
An advanced yet simple attack/defence CTF infrastructure in Go

## Setup Instructions

1.  Enable minikube ingress addon.
   
    ```bash
    $ minikube addons enable ingress
    ```

2. Expose ingress controller.

    ```bash 
    $ kubectl apply -f manifests/dev/expose-controller.yml
    ```

3. Get ip for ingress controller.

    ```bash
    $ minikube service nginx-ingress-controller --url -n kube-system 
    ```

4. Using the ip from step3 and edit `/etc/hosts` and setup local DNS as `challengedeployer.katana.local`. Somewhat similar to this `192.168.49.2    challengedeployer.katana.local`.

5. Now, Katana cluster is up and running and manifests could be applied from `manifests/templates`. 
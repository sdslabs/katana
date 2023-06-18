---
title: "Design"
resources:
  - name: arch
    src: "../../resources/_gen/images/arch.svg"
    title: Architecture
---

Katana uses a namespace-per-team model. Each team is assigned a namespace, and all of the team's resources are deployed into that namespace. This model allows Katana to provide a secure environment for each team, while also allowing teams to interact with each other.

Every team starts with a team pod. The team pod is a pod that is deployed into the team's namespace. The team pod is used to provide the team with a persistent environment. The team pod is also used to provide the team with a persistent storage volume. The team pod is deployed using a [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/), which ensures that the pod is always deployed to the same node. This ensures that the team's persistent storage volume is always available to the team.

The teams are give SSH access to the team pod. Each teams is given a user-password pair that can be used to SSH into the team pod. The team pod is given a public IP address, which can be used to SSH into the team pod from outside of the cluster.

Challenges are pods that are deployed into the team's namespace. On patching, the pod is redeployed.

Katana has its own namespace. This namespace is used to deploy Katana components. These components include flag handler service, challenge checker service, logging service, git server. We will discuss these components in more detail in the next section.

![Image not found](/arch.png)
---
title: "KatanaD"
---

Each team under it's own namespace has a master pod. We have names the general structure of a master pod as KatanaD. Each team will have access to it's own master pod which will contain important directories

1. Challenge Files
2. Another item //Update these after diagram.
3. Another item //Update these after diagram.
4. And another item. //Update these after diagram.

//UPDATE DIAGRAM OF A MASTER PODS WITH THINGS INSIDE IT.
![Image Not Found](/team-pods-architecture.png)

## Structure

The DockerFile can be accessed from [here.](https://github.com/sdslabs/katanad/blob/mainpod/Dockerfile)

## General Flow

- Challenge gets copied to master pod (VANSH)
- [Listening Service unzips the challenge in the master pod (VANSH) ](../KatanaD/v1.md)

Now let's look how a challenge pod gets deployed from the master pod.

- [Challenge-Containers v1](/KatanaD/v1/)

- [Challenge-Containers v2](/KatanaD/v2)

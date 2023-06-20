---
title: "Setup"
---

Setting up Katana involves of the following steps: 

{{<  toc  >}}

### Infraset 

This is setting up the basic pods of katana in katana namespace. These include
- #### MYSQL
  MYSQL picks up the admin config from Config.Toml and is exposed using a NodePort type service on port 32001. This is used to store data for Gogs 
- #### MongoDB
  MongoDB is also a NodePort type service exposed on port 32000. This is used to store team credentials for [SSH Service](/Services/ssh/). This will also be used to store the flags and points of each team. 
- #### GOGS
  This is the locally running GitServer which is used for Patching Service. This has its database in MYSQL and is exposed using Cluster IP. 

  
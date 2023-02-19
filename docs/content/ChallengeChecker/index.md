---
title: "Challenge Checker"
---

# WIP

The challenge checker will be responsible for running the checks against the challenges. The challenge checker will be deployed as a Kubernetes CronJob/Service. The CronJob/Service will run at every tick and will check the status of the challenges. The challenge checker will be responsible for checking the status of the challenges and updating the status of the challenges in the database.

It has been decided to use a Pod in the master namespace to routinely run a knative service which would start a new instantaneous pod and run the checks. The instantaneous pod will return a success or a failure response for its respective request. There will be (no. of challenges x no. of teams) instantaneous pods possible at any given time.

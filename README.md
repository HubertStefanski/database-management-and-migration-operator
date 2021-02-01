# Database Management And Migration Operator

A simple Kubernetes Operator based on the Operator-SDK to manage and migrate databases between cloud environments.


## Requirements
- Running ***Kubernetes/Openshift*** cluster
- Operator-SDK version ***v1.3.0***
- Go version ***go1.15.2***

## Container/Operator Image

The Operator Image can be found at :

 ***hubertstefanski/dbmmo:latest*** (hosted on hub.docker.com)

## (Warning)

This Operator is heavily ***Work In Progress*** as part of a final year project for my bachelors computing degree.
(I cannot accept outside help for the near future). However, feel free to fork/clone this repo to your hearts content.

## Overview

![Operator Overview](documentation/images/operator-overview.png)

## Current functionality
- Hard coded (non-configurable) installation of MySQL Service, Deployment and PersistentVolumeClaim (Update & Delete TBA)
 - Can be deployed

## Planned functionality

- Local database management and provisioning for MySQL(on-cluster)
- Cloud database management and provisioning for Azure (to be expanded to GCP/AWS/etc in the future)
- Data migration between environments (local-Azure / Azure-local etc.)
- Expanded configuration for deployments (non-priority)
- Table management

## Running the Operator

### Local

Currently able to create a deployment, mysql pod, service and privatevolume claim

Prepare the cluster by generating all manifests/code and applying CRDs to the cluster:

 ```bash
 make install
 ```
 
Apply the dbmmo_mysql resource with `oc apply -f example/mysql/dbmm_mysql.yaml -n <NAMESPACE>`


Run the operator locally with

run:

```bash
make run
```

### Deployment

***TBA***


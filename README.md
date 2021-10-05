# Database Management And Migration Operator

# USE AT YOUR OWN RISK 

A simple Kubernetes Operator based on the Operator-SDK to manage and migrate databases between cloud environments.

## Requirements

- Running ***Kubernetes/Openshift*** cluster
- Operator-SDK version ***v1.3.0***
- Go version ***go1.15.2***

## Container/Operator Image

The Operator Image can be found at :

***hubertstefanski/dbmmo:latest*** (hosted on hub.docker.com)

## (Warning)

This Operator is heavily ***Work In Progress*** as part of a final year project for my bachelors computing degree. Use at your own risk, but feel free to contribute features, bug fixes or fork it for your own use!
## Overview

![Operator Overview](documentation/images/operator-overview.png)

## Current functionality

- Hard coded (non-configurable) installation of MySQL Service, Deployment and PersistentVolumeClaim (Update & Delete
  TBA)
- Local database management and provisioning for MySQL(on-cluster) accessible through service
- Cloud database creation and deletion for Azure (partial configuration)
- Database migration
    - Azure -> OnCluster
    - OnCluster -> Azure

## Planned functionality
- Data migration between environments (OnCluster -> Azure / Azure -> OnCluster)
- Expanded configuration for deployments (non-priority)
- Table management

## Running the Operator

### Local

Currently able to create a Deployment, MySQL Pod, Service and Persistent Volume Claim

Prepare the cluster by generating all manifests/code and applying CRDs to the cluster:

 ```bash
 make install
 ```

Apply the dbmmo_mysql resource with

```bash
kubectl apply -f example/mysql/dbmm_mysql.yaml -n <NAMESPACE>
```

Run the operator locally with

run:

```bash
make run
```

### Deployment

Prepare the cluster by running

```bash
make cluster/prepare/local NAMESPACE=<NAMESPACE>
 ```

This will create all necessary roles, service accounts, roles and bindings for the operator to be able to run as a
deployment

Apply the operator deployment to the cluster with:

***NOTE: Ensure that the operator is being deployed to the same namespace that was passed in the previous env var***

```bash
kubectl apply -f example/operator/operator_<version or latest>.yaml -n <NAMESPACE>
```


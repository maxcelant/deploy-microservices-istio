# Istio Microservice Kubernetes Sample

Quickly deploys three microservices to different environments of your choice. 

- `/src` holds the each microservice's specific code
- `/deploy-docker` deploying these services to a Docker network
- `/scripts` holds any scripts related to building, running, etc

The goal of this project is to deploy these services to a KinD cluster and tune their interactions using Istio's service mesh capabilities. 
# Deploy Microservices in Kubernetes

Quickly deploys three microservices to different environments of your choice. 

- `/src` holds the each microservice's specific code. They each run their own dependent postgres database. 
- `/deploy-docker` deploying these services to a Docker network
- `/scripts` holds any scripts related to building, running, etc

The goal of this project is to deploy these services to a KinD cluster and tune their interactions using Istio's service mesh capabilities. 

### Running Locally
This is if you want to run each microservice locally on a different port. There are some pre-requisites to do so. 

1. Make sure you have Go and Docker Desktop installed on your local machine.

2. Run the scripts `./scripts/docker-build-dbs` followed by `./scripts/docker-migrate-dbs` to initialize your postgres databases.

3. Using your preferred terminal window manager, start each service in their own window by running `make run`. They are preset to listen on separate ports (8080, 8081, 8082).

4. Make REST calls from your preferred client.
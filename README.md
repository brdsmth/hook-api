## Introduction

The `hook-api` is an HTTP web server that allows users to add jobs to the queue

### Usage

The `hook-api` is meant to be run in conjunction with the `hook-scheduler` and `hook-runner` services.

To run the `hook-api` server the following `.env` variables need to be set

```
AWS_CONFIG_PROFILE=
DYNAMODB_QUEUE_TABLE=
```

Once these are added the server can be started by

```
go run main.go
```

Once the `hook-api` server is running you can send a **POST** request to the `/add` endpoint which will add the job to the job queue in **DynamoDB**

## Deployment

### Kubernetes

#### Local

Start `minikube`

```
minikube start
```

Direct `minikube` to use the `docker` env. Any `docker build ...` commands after this command is run will build inside the `minikube` registry and will not be visible in Docker Desktop. `minikube` uses its own docker daemon which is separate from the docker daemon on your host machine. Running `docker images` inside the `minikube` vm will show the images accessible to `minikube`

```
eval $(minikube docker-env)
```

```
docker build -t hook-api-image:latest .
```

#### Environment Variables (if needed)

```
kubectl create secret generic awsconfig-secret --from-env-file=./.env
kubectl create secret generic dynamodbqueuetable-secret --from-env-file=./.env
```

```
kubectl apply -f ./k8s/hook-api.deployment.yaml
```

```
kubectl apply -f ./k8s/hook-api.service.yaml
```

```
kubectl get deployments
```

```
kubectl get pods
```

```
minikube service hook-api-service
```

After running the last comment the application will be able to be accessed in the browser at the specified port that `minikube` assigns.

#### Troubleshooting

```
minikube ssh 'docker images'
```

```
kubectl logs <pod-name>
```

```
kubectl logs -f <pod-name>
```

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

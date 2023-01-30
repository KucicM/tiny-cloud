
# tiny-cloud

Run docker task in a cloud

## Setup

Create user "tiny-cloud" in AWS console. IAM->USERS

Add permissions:

- AmazonEC2ContainerRegistryFullAccess

## Run

Build docker image:

``` docker
docker -t hello-world .
```

Run a task:

``` bash
go run . --image hello-world
```

## Delete resources

``` bash
go run . --destroy true
```

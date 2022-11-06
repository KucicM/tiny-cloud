
# tiny-cloud

CLI wrapper to run simple tasks on a cloud

## Install

prerequisites:
Docker:

```text
https://docs.docker.com/engine/install/
```

AWS CLI:

```text
https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
```

Terraform:

```text
https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli
```

## Setup

Create user "tiny-cloud" in AWS console. IAM->USERS

Add permissions:

- AmazonEC2ContainerRegistryFullAccess

Run:

``` bash
aws configure --profile tiny-cloud
```

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

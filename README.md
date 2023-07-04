# tiny-cloud

Deploy docker image on remote machine like AWS EC2.


NOTE: Work in progress!

## Setup

### AWS

Create user "tiny-cloud" in AWS console. IAM->USERS

Create role "tiny-cloud" in AWS console. IAM->ROLES

Add permissions:

- AmazonEC2ContainerRegistryFullAccess

### Configure

``` bash
tiny-cloud configure
```

## Useage

Run

``` bash
tiny-cloud run <image-name> <override options>
```

## Delete resources

``` bash
tiny-cloud --destroy true
```

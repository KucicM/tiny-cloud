
# tiny-cloud

Run task on cloud VMs.

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

## Run

Run a task:

``` bash
tiny-cloud run <command>
```

## Delete resources

``` bash
tiny-cloud --destroy true
```

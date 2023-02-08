# Hello world example

## Usage with AWS

``` bash
tiny-cloud run --src-path=hello_world.sh --vm-type=t2.micro --data-out=hello.txt
```

## Explained

1. Start `t2.micro` VM on AWS
2. Copy `hello_world.sh` file to the VM
3. Execute on the VM
4. Copy hello.txt file from VM to the S3 bucket

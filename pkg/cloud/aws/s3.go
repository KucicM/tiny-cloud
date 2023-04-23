package aws

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketPrefix string = "tiny-cloud-"
const maxSizeName int = 32

func CreateS3(req AwsSetupRequest) error {
	creds := credentials.NewStaticCredentialsProvider(
		req.AccessKeyId,
		req.SeacretAccessKey,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(req.Region),
		config.WithCredentialsProvider(creds),
	)

    if err != nil {
        return err
    }

    client := s3.NewFromConfig(cfg)


    res, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    if err != nil {
        return err
    }

    for _, bucket := range res.Buckets {
        if strings.HasPrefix(*bucket.Name, bucketPrefix) {
            return nil
        }
    }

    bucketName := createRandomName()
    out, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
        CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(req.Region),
		},
    })

    fmt.Printf("%v\n", out)
    return err
}

func createRandomName() string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, maxSizeName - len(bucketPrefix))
    rand.Read(b)
    return fmt.Sprintf("%s%x", bucketPrefix, b)
}

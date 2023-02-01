package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type AWS struct {
	ops *tinycloud.Ops
	cfg aws.Config
}

func New() *AWS {
	// maybe add cofigurable profile?

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(tinycloud.PROFILE_NAME),
	)

	if err != nil {
		log.Fatalln(err)
	}

	return &AWS{cfg: cfg}
}

func randomSuffix(n int) string {
	rand.Seed(time.Now().Unix())
	out := make([]byte, n)
	for i := range out {
		out[i] = byte(97 + rand.Intn(26))
	}
	return string(out)
}

func getIp() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bs)) + "/32", nil
}

func (a *AWS) Run(ops tinycloud.Ops) error {
	ip, err := getIp()
	if err != nil {
		return err
	}
	// prepare docker image

	suffix := randomSuffix(8)
	name := fmt.Sprintf("%s-%s", tinycloud.PROFILE_NAME, suffix)

	ec2Client := ec2.NewFromConfig(a.cfg)

	// create new security group
	sgOps := &ec2.CreateSecurityGroupInput{
		Description: aws.String("tiny-cloud"),
		GroupName:   aws.String(name),
	}

	_, err = ec2Client.CreateSecurityGroup(context.TODO(), sgOps, func(o *ec2.Options) {})
	if err != nil {
		log.Printf("error %s", err)
		return err
	}

	// set security group rules (ssh)
	ingressOps := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(name),
		IpPermissions: []types.IpPermission{{
			FromPort:   aws.Int32(22),
			ToPort:     aws.Int32(22),
			IpProtocol: aws.String("TCP"),
			IpRanges: []types.IpRange{{
				CidrIp: aws.String(ip),
			}},
		}},
	}

	_, err = ec2Client.AuthorizeSecurityGroupIngress(context.TODO(), ingressOps, func(o *ec2.Options) {})
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	// vmReq := EC2Request{InstanceType: ops.VmType}

	// securityGroupName := "tiny-cloud"
	// op := &ec2.DescribeSecurityGroupsInput{
	// 	Filters: []types.Filter{{
	// 		Name:   aws.String("group-name"),
	// 		Values: []string{securityGroupName},
	// 	}},
	// }

	// out, err := ec2Client.DescribeSecurityGroups(context.TODO(), op, func(o *ec2.Options) {})
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// log.Printf("%+v", out.SecurityGroups)

	// op := &ec2.DescribeKeyPairsInput{KeyNames: []string{tinycloud.PROFILE_NAME}}
	// out, err := ec2Client.DescribeKeyPairs(context.TODO(), op, func(o *ec2.Options) {})
	// if err != nil {
	// 	return err
	// }

	// if len(out.KeyPairs) == 0 {
	// 	// todo create key
	// 	return nil
	// }

	// info := out.KeyPairs[0]
	// name := info.KeyName

	// log.Printf("out %v\n", len(out.KeyPairs))
	// log.Printf("out %v\n", *name)

	// ec2Client.DescribeKeyPairs(ec2Client)

	// key, err := NewKeyPair(ec2Client)

	NewVm(ec2Client, VmRequest{ops.VmType, true, "test-keys", name})

	// todo push docker image

	// clean docker staging image

	return nil
}

func (a *AWS) Destroy() error {
	// DestroyVMs(a.cfg)
	// keys
	return nil
}

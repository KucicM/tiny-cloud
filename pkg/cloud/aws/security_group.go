package aws

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroupApi interface {
	CreateSecurityGroup(ctx context.Context,
		params *ec2.CreateSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
	DeleteSecurityGroup(ctx context.Context,
		params *ec2.DeleteSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
}

func CreateSecurityGroup(groupName string, api SecurityGroupApi) error {
	ops := &ec2.CreateSecurityGroupInput{
		Description: aws.String("tiny-cloud"),
		GroupName:   aws.String(groupName),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeSecurityGroup,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(groupName),
			}},
		}},
	}

	_, err := api.CreateSecurityGroup(context.TODO(), ops, opsFn)
	return err
}

func DeleteSecurityGroups(groupNames []string, api SecurityGroupApi) error {
	for _, groupName := range groupNames {
		ops := &ec2.DeleteSecurityGroupInput{GroupName: aws.String(groupName)}
		_, err := api.DeleteSecurityGroup(context.TODO(), ops, opsFn)
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			return err
		}
	}
	return nil
}

type AuthorizeSecurityGroupApi interface {
	AuthorizeSecurityGroupIngress(ctx context.Context,
		params *ec2.AuthorizeSecurityGroupIngressInput,
		optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}

func AuthorizeSSH(groupName string, api AuthorizeSecurityGroupApi) error {
	ip, err := getIp()
	if err != nil {
		return err
	}

	ops := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(groupName),
		IpPermissions: []types.IpPermission{{
			FromPort:   aws.Int32(22),
			ToPort:     aws.Int32(22),
			IpProtocol: aws.String("TCP"),
			IpRanges: []types.IpRange{{
				CidrIp: aws.String(ip),
			}},
		}},
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeSecurityGroupRule,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(groupName),
			}},
		}},
	}

	_, err = api.AuthorizeSecurityGroupIngress(context.TODO(), ops, opsFn)
	return err
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

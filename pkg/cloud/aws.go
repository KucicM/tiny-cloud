package cloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

type Ec2 struct {
	Id      string
	SSHKey  []byte
	DNSName string
}

func (e *Ec2) Run(task tinycloud.TaskDefinition) {

}

type AwsSetupRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	InstanceType     string
	Iam              string
}

// CERATE AWS VM

func StartAwsVm(req AwsSetupRequest) (tinycloud.Vm, error) {
	// auth
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
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	runId, err := data.GetNewRunId(req.ProfileName)
	if err != nil {
		return nil, err
	}

	if err = CreateSecurityGroup(runId, client); err != nil {
		return nil, err
	}

	if err = AuthorizeSSH(runId, client); err != nil {
		return nil, err
	}

	var sshKey []byte
	if sshKey, err = GenerateSSHKey(runId, client); err != nil {
		return nil, err
	}
	data.AddPemKey(runId, sshKey)

	var instanceId string
	if instanceId, err = CretaeInstance(runId, req.InstanceType, req.Iam, client); err != nil {
		return nil, err
	}

	err = waitInstanceStart(instanceId, client)
	if err != nil {
		return nil, err
	}

	var dnsName string
	if dnsName, err = getDNSName(instanceId, client); err != nil {
		return nil, err
	}

	return &Ec2{
		Id:      runId,
		SSHKey:  sshKey,
		DNSName: dnsName,
	}, nil
}

type SecurityGroupApi interface {
	CreateSecurityGroup(ctx context.Context,
		params *ec2.CreateSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
	DeleteSecurityGroup(ctx context.Context,
		params *ec2.DeleteSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
}

func CreateSecurityGroup(runId string, api SecurityGroupApi) error {
	ops := &ec2.CreateSecurityGroupInput{
		Description: aws.String("tiny-cloud"),
		GroupName:   aws.String(runId),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeSecurityGroup,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(runId),
			}},
		}},
	}

	_, err := api.CreateSecurityGroup(context.TODO(), ops, opsFn)
	return err
}

type AuthorizeSecurityGroupApi interface {
	AuthorizeSecurityGroupIngress(ctx context.Context,
		params *ec2.AuthorizeSecurityGroupIngressInput,
		optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}

func AuthorizeSSH(runId string, api AuthorizeSecurityGroupApi) error {
	ip, err := getIp()
	if err != nil {
		return err
	}

	ops := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(runId),
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
				Value: aws.String(runId),
			}},
		}},
	}

	_, err = api.AuthorizeSecurityGroupIngress(context.TODO(), ops, opsFn)
	return err
}

type SSHKeyApi interface {
	CreateKeyPair(ctx context.Context,
		params *ec2.CreateKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(ctx context.Context,
		params *ec2.DeleteKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteKeyPairOutput, error)
}

func GenerateSSHKey(runId string, api SSHKeyApi) ([]byte, error) {
	ops := &ec2.CreateKeyPairInput{
		KeyName: aws.String(runId),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeKeyPair,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(runId),
			}},
		}},
	}
	out, err := api.CreateKeyPair(context.TODO(), ops, opsFn)
	if err != nil {
		return nil, err
	}

	return []byte(*out.KeyMaterial), nil
}

type InstanceApi interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	DescribeInstanceStatus(ctx context.Context,
		params *ec2.DescribeInstanceStatusInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)

	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func CretaeInstance(runId, instanceType, iam string, api InstanceApi) (string, error) {
	ops := &ec2.RunInstancesInput{
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		ImageId:                           aws.String(iam), // TODO
		InstanceType:                      types.InstanceType(instanceType),
		KeyName:                           aws.String(runId),
		SecurityGroups:                    []string{runId},
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(runId),
			}},
		}},
	}

	out, err := api.RunInstances(context.TODO(), ops, opsFn)
	if err != nil {
		return "", fmt.Errorf("faild to create new vm with error: %s", err)
	}

	if out == nil {
		return "", fmt.Errorf("run instances return nil")
	}

	if len(out.Instances) == 0 {
		return "", fmt.Errorf("run instances return zero instances")
	}

	instanceId := *out.Instances[0].InstanceId
	if instanceId == "" {
		return "", fmt.Errorf("run instaces return instance without id")
	}

	return instanceId, nil
}

func waitInstanceStart(instanceId string, api InstanceApi) error {
	fmt.Println("waiting instance start")
	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         []string{instanceId},
	}

	// TODO add timeout?
	for {
		out, err := api.DescribeInstanceStatus(context.TODO(), ops, opsFn)

		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		if out == nil || len(out.InstanceStatuses) == 0 {
			continue
		}

		state := out.InstanceStatuses[0].InstanceState
		if state.Name == types.InstanceStateNameRunning {
			break
		}
	}

	return nil
}

func getDNSName(instanceId string, api InstanceApi) (string, error) {
	ops := &ec2.DescribeInstancesInput{InstanceIds: []string{instanceId}}

	out, err := api.DescribeInstances(context.TODO(), ops, opsFn)
	if err != nil {
		return "", fmt.Errorf("cannot get dns name of instanceId: %s, error: %s", instanceId, err)
	}

	if out == nil || len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("got unexpected result from describe instances %+v", out)
	}

	instance := out.Reservations[0].Instances[0]
	return *instance.PublicDnsName, nil
}

// DELETE AWS VM
type AwsDestroyRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	RunIds           []string
}

// deletes all resoures created by the user
func DestroyAws(req AwsDestroyRequest) error {
	// auth
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

	client := ec2.NewFromConfig(cfg)

	if err = deleteInstances(client, req.RunIds); err != nil {
		return err
	}

	if err = deleteSSHKeys(client, req.RunIds); err != nil {
		return err
	}

	if err = deleteSecurityGroups(client, req.RunIds); err != nil {
		return err
	}

	return nil
}

func deleteInstances(api InstanceApi, runIds []string) error {
	findOps := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("tag-key"), Values: []string{"tiny-cloud"}},
			{Name: aws.String("instance-state-name"), Values: []string{"pending", "running", "stopped"}},
		},
	}

	des, err := api.DescribeInstances(context.TODO(), findOps, opsFn)
	if err != nil {
		return err
	}

	if des == nil || len(des.Reservations) == 0 {
		return nil
	}

	instanceIds := make([]string, 0)
	for _, reservation := range des.Reservations {
		for _, instance := range reservation.Instances {
			if hasRunIdTag(instance.Tags, runIds) {
				instanceIds = append(instanceIds, *instance.InstanceId)
			}
		}
	}

	ops := &ec2.TerminateInstancesInput{InstanceIds: instanceIds}
	if _, err = api.TerminateInstances(context.TODO(), ops, opsFn); err != nil {
		return err
	}

	return waitInstanceTermination(api, instanceIds)
}

func waitInstanceTermination(api InstanceApi, instanceIds []string) error {
	fmt.Println("wating instance termination")
	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         instanceIds,
	}

	// todo timeout?
	anyUp := true
	for anyUp {
		out, err := api.DescribeInstanceStatus(context.TODO(), ops, opsFn)

		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		if out == nil {
			continue
		}

		anyUp = false
		for _, instance := range out.InstanceStatuses {
			if instance.InstanceState.Name != types.InstanceStateNameTerminated {
				anyUp = true
			}
		}
	}
	return nil
}

func hasRunIdTag(tags []types.Tag, runIds []string) bool {
	for _, tag := range tags {
		if *tag.Key == "tiny-cloud" {
			for _, runId := range runIds {
				if runId == *tag.Value {
					return true
				}
			}
		}
	}
	return false
}

func deleteSSHKeys(api SSHKeyApi, runIds []string) error {
	// runIds == keyNames
	for _, runId := range runIds {
		ops := &ec2.DeleteKeyPairInput{KeyName: aws.String(runId)}
		_, err := api.DeleteKeyPair(context.TODO(), ops, opsFn)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteSecurityGroups(api SecurityGroupApi, runIds []string) error {
	// runIds == securityGroupName
	for _, runId := range runIds {
		ops := &ec2.DeleteSecurityGroupInput{GroupName: aws.String(runId)}
		_, err := api.DeleteSecurityGroup(context.TODO(), ops, opsFn)
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			return err
		}
	}
	return nil
}

// func (a *AWS) Run(ops tinycloud.Ops) error {

// 	// connect to vm and execute command

// 	pemBytes := []byte(*keyOut.KeyMaterial)
// 	signer, err := ssh.ParsePrivateKey(pemBytes)
// 	if err != nil {
// 		log.Fatalf("parse key failed:%v", err)
// 	}
// 	config := &ssh.ClientConfig{
// 		User:            "ec2-user",
// 		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	url := fmt.Sprintf("%s:22", dnsName)
// 	var conn *ssh.Client
// 	for {
// 		conn, err = ssh.Dial("tcp", url, config)
// 		if err == nil {
// 			break
// 		}
// 		time.Sleep(time.Second)
// 	}
// 	defer conn.Close()

// 	session, err := conn.NewSession()
// 	if err != nil {
// 		log.Fatalf("session failed:%v", err)
// 	}
// 	defer session.Close()

// 	stdin, err := session.StdinPipe()
// 	if err != nil {
// 		return err
// 	}
// 	defer stdin.Close()

// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr

// 	if err := session.Shell(); err != nil {
// 		return err
// 	}

// 	// push to S3?

// 	stdin.Write([]byte("whoami\n"))
// 	// stdin.Write([]byte("sudo shutdown now\n"))
// 	session.Wait()
// 	return nil
// }

func opsFn(*ec2.Options) {
	// empty
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

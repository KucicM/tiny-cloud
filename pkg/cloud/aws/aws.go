package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

// CERATE AWS VM

type AwsSetupRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	InstanceType     string
	Iam              string
}

func StartVm(req AwsSetupRequest) (tinycloud.Vm, error) {
	// auth
	client, err := getClient(req.Region, req.AccessKeyId, req.SeacretAccessKey)
	if err != nil {
		return nil, err
	}

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
	if sshKey, err = CreateKeyPair(runId, client); err != nil {
		return nil, err
	}
	data.AddPemKey(runId, sshKey)

	var instanceId string
	if instanceId, err = CretaeEc2(runId, req.InstanceType, req.Iam, client); err != nil {
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

// DELETE AWS VM
type AwsDestroyRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	RunIds           []string
}

// deletes resoures created by the user
func DestroyAws(req AwsDestroyRequest) error {
	// auth
	client, err := getClient(req.Region, req.AccessKeyId, req.SeacretAccessKey)
	if err != nil {
		return err
	}

	if err = DeleteEc2(req.RunIds, client); err != nil {
		return err
	}

	if err = DeleteKeyPairs(req.RunIds, client); err != nil {
		return err
	}

	if err = DeleteSecurityGroups(req.RunIds, client); err != nil {
		return err
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

func getClient(region, accessKeyId, seacretAccessKey string) (*ec2.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		accessKeyId,
		seacretAccessKey,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	return client, err
}

func opsFn(*ec2.Options) {
	// empty
}

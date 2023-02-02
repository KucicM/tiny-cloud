package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"golang.org/x/crypto/ssh"
)

type AWS struct {
	ops *tinycloud.Ops
	cfg aws.Config
}

const PROFILE_NAME = "tiny-cloud"

func New() *AWS {
	// maybe add cofigurable profile?

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(PROFILE_NAME),
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
	name := fmt.Sprintf("%s-%s", PROFILE_NAME, suffix)

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

	// keys
	keyOps := &ec2.CreateKeyPairInput{KeyName: aws.String(name)}
	keyOut, err := ec2Client.CreateKeyPair(context.TODO(), keyOps, func(o *ec2.Options) {})
	if err != nil {
		return err
	}

	// create vm
	vmOps := &ec2.RunInstancesInput{
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		ImageId:                           aws.String("ami-06c39ed6b42908a36"), // TODO
		InstanceType:                      types.InstanceType(ops.VmType),
		KeyName:                           aws.String(name),
		SecurityGroups:                    []string{name},
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
	}
	out, err := ec2Client.RunInstances(context.TODO(), vmOps, func(o *ec2.Options) {})
	if err != nil {
		return fmt.Errorf("faild to create new vm with error: %s", err)
	}

	if len(out.Instances) == 0 {
		return fmt.Errorf("got zero vm after create")
	}

	vmId := *out.Instances[0].InstanceId
	log.Println("vmId:", vmId)

	// wait to start vm
	statusOps := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         []string{vmId},
	}

	log.Println("wait vm start")
	for {
		statusOut, err := ec2Client.DescribeInstanceStatus(context.TODO(), statusOps, func(o *ec2.Options) {})
		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		// is started?
		if *statusOut.InstanceStatuses[0].InstanceState.Code == 16 {
			break
		}
	}

	// get dns name
	dnsOps := &ec2.DescribeInstancesInput{InstanceIds: []string{vmId}}

	dnsOut, err := ec2Client.DescribeInstances(context.TODO(), dnsOps, func(o *ec2.Options) {})
	if err != nil {
		return err
	}

	dnsName := *dnsOut.Reservations[0].Instances[0].PublicDnsName
	log.Println("dns name:", dnsName)

	// connect to vm and execute command

	pemBytes := []byte(*keyOut.KeyMaterial)
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig{
		User:            "ec2-user",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	url := fmt.Sprintf("%s:22", dnsName)
	var conn *ssh.Client
	for {
		conn, err = ssh.Dial("tcp", url, config)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("session failed:%v", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		return err
	}

	stdin.Write([]byte("whoami\n"))
	stdin.Write([]byte("sudo shutdown now\n"))
	session.Wait()
	return nil
}

func (a *AWS) Destroy() error {
	// DestroyVMs(a.cfg)
	// keys
	return nil
}

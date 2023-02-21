package cloud

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"golang.org/x/crypto/ssh"
)

// check if docker image exists (do not create any resources if it does not exist)
// save docker image into a "file"
// push docker "file" in vm
// unpack docker image
// start docker container
// wait to finish?

func Run(task tinycloud.TaskDefinition) error {
	signer, err := ssh.ParsePrivateKey(task.SSHKey)
	if err != nil {
		return err
	}

	cfg := &ssh.ClientConfig{
		User:            "ec2-user",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	url := fmt.Sprintf("%s:22", task.DNSName)

	var conn *ssh.Client
	for {
		if conn, err = ssh.Dial("tcp", url, cfg); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	defer conn.Close()

	client, err := scp.NewClientBySSH(conn)
	if err != nil {
		return err
	}

	img, err := GetDockerImage(task.DockerImageId)
	if err != nil {
		return err
	}
	defer img.Close()

	imageName := fmt.Sprintf("tiny-cloud-docker-image-%s", task.DockerImageId)
	err = client.CopyFilePassThru(
		context.Background(),
		img,
		imageName,
		"0655",
		func(r io.Reader, total int64) io.Reader { return r })

	if err != nil {
		return err
	}
	client.Close()

	script := fmt.Sprintf("docker load --input %s\n", imageName)
	script += fmt.Sprintf("docker run %s\n", task.DockerImageId)
	script += "exit\n"

	// run docker image on vm
	log.Println("create session")
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	if err := session.Shell(); err != nil {
		return err
	}

	stdin.Write([]byte(script))
	session.Wait()

	return nil
}

func GetDockerImage(imageId string) (io.ReadCloser, error) {
	if err := os.Setenv(client.EnvOverrideAPIVersion, "1.41"); err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	reader, err := cli.ImageSave(context.Background(), []string{imageId})
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func DoseDockerImageExists(imageId string) (bool, error) {
	if err := os.Setenv(client.EnvOverrideAPIVersion, "1.41"); err != nil {
		return false, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, err
	}

	// todo add filter for imageId
	images, err := cli.ImageList(context.TODO(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, img := range images {
		id := strings.Split(img.ID, ":")[1]
		if strings.HasPrefix(id, imageId) {
			return true, nil
		}
	}
	return false, nil
}

package cloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"golang.org/x/crypto/ssh"
)

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

	f, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	var conn *ssh.Client
	for {
		if conn, err = ssh.Dial("tcp", url, cfg); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	defer conn.Close()

	// copy docker image to vm
	client, err := scp.NewClientBySSH(conn)
	if err != nil {
		return err
	}
	if err = client.CopyFromFile(context.Background(), *f, "test.txt", "0655"); err != nil {
		return err
	}
	client.Close()

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

	stdin.Write([]byte("ls\n"))
	session.Wait()

	return nil
}

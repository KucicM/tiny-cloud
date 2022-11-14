package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type ContainerRepository interface {
	Push(imgName string)
	Destroy()
}

const name = "tiny-cloud-repository"

type ECR struct {
	client *ecr.Client
}

func NewRegistry(cfg aws.Config) *ECR {
	return &ECR{
		client: ecr.NewFromConfig(cfg),
	}
}

func (e *ECR) createRepository() {

	ops := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
		Tags: []types.Tag{{
			Key:   aws.String("tiny-cloud"),
			Value: aws.String("repository"),
		}},
	}

	_, err := e.client.CreateRepository(context.TODO(), ops, func(o *ecr.Options) {})
	if err == nil {
		log.Printf("Created new repository '%s'\n", name)
	} else if !strings.Contains(err.Error(), "RepositoryAlreadyExistsException") {
		log.Fatalf("cannot create repository %s, error: %s", name, err)
	}
}

func (e *ECR) Destroy() {
	log.Printf("Removing repository %s\n", name)
	ops2 := &ecr.DeleteRepositoryInput{
		RepositoryName: aws.String(name),
		Force:          true,
	}
	_, err := e.client.DeleteRepository(context.TODO(), ops2, func(o *ecr.Options) {})
	if err != nil {
		log.Printf("error deleting repoistory '%s' %s\n", name, err)
	}
}

// push image (add image name)
func (e *ECR) Push(img string) {

	// create if needed
	e.createRepository()

	// login to docker
	e.dockerLogin()

	// tag image
	log.Println("Tagging docker image")
	tag := fmt.Sprintf("%s:latest", e.uri())
	run(fmt.Sprintf("docker tag %s %s", img, tag))

	// push
	log.Printf("Push image %s to %s\n", img, name)
	out, _ := run(fmt.Sprintf("docker push %s", tag))
	log.Println(string(out))

}

func (e *ECR) dockerLogin() {
	log.Println("docker login")
	query := "docker login --username AWS --password %s %s"
	out, _ := run(fmt.Sprintf(query, e.token(), e.uri()))
	log.Println(string(out))
}

func (e *ECR) uri() string {
	ops := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{name},
	}
	out, _ := e.client.DescribeRepositories(context.TODO(), ops, func(o *ecr.Options) {})
	for _, r := range out.Repositories {
		return *r.RepositoryUri
	}
	log.Fatalf("Cannot find uri for %s repository\n", name)
	return ""
}

func (e *ECR) token() string {
	ops := &ecr.GetAuthorizationTokenInput{}
	out, err := e.client.GetAuthorizationToken(context.TODO(), ops, func(o *ecr.Options) {})
	if err != nil || len(out.AuthorizationData) == 0 {
		log.Fatalf("Cannot get password for %s %s\n", name, err)
	}
	auth := *out.AuthorizationData[0].AuthorizationToken
	token, err := base64.RawStdEncoding.DecodeString(auth)
	if err != nil {
		log.Fatalf("Error converting auth token from base64 %s\n", err)
	}
	return strings.TrimPrefix(string(token), "AWS:")
}

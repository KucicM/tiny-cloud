package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type ECS struct {
	name   string
	client *ecs.Client
}

func NewECS(cfg aws.Config) *ECS {
	return &ECS{
		name:   fmt.Sprintf("%s-ecs", tinycloud.Name),
		client: ecs.NewFromConfig(cfg),
	}
}

func (e *ECS) create() {
	ops := &ecs.CreateClusterInput{
		ClusterName: aws.String(e.name),
	}

	_, err := e.client.CreateCluster(context.TODO(), ops, func(o *ecs.Options) {})
	if err == nil {
		log.Printf("created new ecs '%s'\n", e.name)
	} else if err != nil {
		log.Fatalf("failed to create cluste %s", err)
	}
}

func (e *ECS) destroy() {
	log.Printf("remove ecs '%s'\n", e.name)

	ops := &ecs.DeleteClusterInput{Cluster: &e.name}
	_, err := e.client.DeleteCluster(context.TODO(), ops, func(o *ecs.Options) {})
	if err != nil {
		log.Fatalf("failed to delete cluster %s\n", err)
	}
}

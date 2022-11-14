package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const clusterName = "tiny-cloud-cluster"

type ElasticContainerService struct {
	client *ecs.Client
}

func NewElasticContainerService(cfg aws.Config) *ElasticContainerService {
	return &ElasticContainerService{
		client: ecs.NewFromConfig(cfg),
	}
}

func (c *ElasticContainerService) Create() {
	ops := &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
		Tags: []types.Tag{{
			Key:   aws.String("tiny-cloud"),
			Value: aws.String("cluster"),
		}},
	}
	_, err := c.client.CreateCluster(context.TODO(), ops, func(o *ecs.Options) {})
	if err != nil {
		log.Fatalf("Failed to create ECS cluster %s\n", err)
	}
}

func (c *ElasticContainerService) Destroy() {
	log.Println("Destroy ECS cluster")
	ops := &ecs.DeleteClusterInput{
		Cluster: aws.String(clusterName),
	}
	_, err := c.client.DeleteCluster(context.TODO(), ops, func(o *ecs.Options) {})
	if err != nil {
		log.Printf("Failed to delete ECS cluster %s\n", err)
	}
}

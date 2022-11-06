package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ContainerRepository interface {
	Push(imgName string)
	DeleteImgs()
}

type descRepoResp struct {
	Repos []repo `json:"repositories"`
}

type repo struct {
	Name string `json:"repositoryName"`
	Uri  string `json:"repositoryUri"`
}

type Images struct {
	Ids []struct {
		Tag string `json:"imageTag"`
	} `json:"ImageIds"`
}

type ECR struct {
}

// push image (add image name)
func (r *ECR) Push(img string) {
	// login to docker
	r.dockerLogin()

	// tag image
	tag := r.nextTag()
	out, _ := run(fmt.Sprintf("docker tag %s %s", img, tag))
	log.Println(string(out))

	// push
	out, _ = run(fmt.Sprintf("docker push %s", tag))
	log.Println(string(out))

}

// delete images
func (r *ECR) DeleteImgs() {
	// get images
	imgs := r.images()

	tags := make([]string, 0)
	for _, img := range imgs.Ids {
		tags = append(tags, fmt.Sprintf("imageTag=%s", img.Tag))
	}

	// bulk delete
	if len(tags) != 0 {
		ts := strings.Join(tags, " ")
		_, _ = run(fmt.Sprintf("aws ecr batch-delete-image --repository-name %s --image-ids %s", repoName, ts))
	}
}

func (r *ECR) nextTag() string {
	return fmt.Sprintf("%s:%d", r.uri(), r.nextVer())
}

func (r *ECR) nextVer() int64 {
	imgs := r.images()
	var maxTag int64 = 0
	for _, d := range imgs.Ids {
		if t, err := strconv.ParseInt(d.Tag, 10, 64); err == nil && t > maxTag {
			maxTag = t
		}
	}
	return maxTag + 1
}

// TODO handle errors
func (r *ECR) images() Images {
	out, _ := run(fmt.Sprintf("aws ecr list-images --repository-name %s --profile %s", repoName, profile))
	var res Images
	_ = json.Unmarshal(out, &res)
	return res
}

func (r *ECR) dockerLogin() {
	log.Println("docker login")
	query := "aws ecr get-login-password --profile %s | docker login --username AWS --password-stdin %s"
	out, _ := run(fmt.Sprintf(query, profile, r.uri()))
	log.Println(string(out))
}

func (r *ECR) uri() string {
	return r.tinyRepo().Uri
}

func (r *ECR) tinyRepo() repo {
	repos := r.fetchRepos()
	for _, re := range repos.Repos {
		if re.Name == repoName {
			return re
		}
	}
	return repo{}
}

func (r *ECR) fetchRepos() descRepoResp {
	out, _ := run(fmt.Sprintf("aws ecr describe-repositories --profile %s", profile)) // todo chec kerror

	var res descRepoResp
	_ = json.Unmarshal(out, &res) // todo check error
	return res
}

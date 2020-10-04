package service

import (
	"fmt"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/utils"
)

const (
	Up   action = "UP"
	Down action = "DOWN"
)

func makeDockerDeploymentJob(config config.Config, action action) (job.Job, error) {
	switch action {
	case Up:
		return createContainersJob(config), nil
	case Down:
		return deleteContainersJob(config), nil
	default:
		return job.Job{}, fmt.Errorf("unknown action %s", action)
	}
}

func createContainersJob(config config.Config) job.Job {
	jobsArgs := job.Args{"config": config}
	return job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   config.Name,
		Type:        ContainerCreate,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         createContainerResources,
	}
}

func deleteContainersJob(config config.Config) job.Job {
	jobsArgs := job.Args{"namespace": config.Name}
	return job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   config.Name,
		Type:        ContainerDelete,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         deleteContainerResources,
	}
}
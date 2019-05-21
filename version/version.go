package version

import (
	"errors"
	"os"
)

type GetVersion interface {
	GetDeploymentVersion(envname string) (*DeploymentVersion, error)
}

type DeploymentVersion struct {
	commitID string `json:"commit_id"`
}

func (d *DeploymentVersion) GetDeploymentVersion(envname string) (*DeploymentVersion, error) {
	dep := &DeploymentVersion{}
	v := os.Getenv(envname)

	if v == "" {

		return nil, errors.New("env is not set")

	}

	dep.commitID = v

	return dep, nil

}

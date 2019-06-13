package version

import (
	"errors"
	"os"
)

type GetVersion interface {
	GetDeploymentVersion(envname string) (*DeploymentVersion, error)
}

func NewDeploymentVersion() GetVersion {

	return &DeploymentVersion{}

}

type DeploymentVersion struct {
	CommitID string `json:"commit_id"`
}

func (d *DeploymentVersion) GetDeploymentVersion(envname string) (*DeploymentVersion, error) {
	v := os.Getenv(envname)

	if v == "" {

		return nil, errors.New("commit id is not set")

	}

	d.CommitID = v

	return d, nil

}

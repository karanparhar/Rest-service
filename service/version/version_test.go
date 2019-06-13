package version

import (
	"os"
	"testing"
)

func TestVersion_GetDeploymentVersion(t *testing.T) {

	var version = DeploymentVersion{}

	os.Setenv("GitCommit", "testing")

	out, err := version.GetDeploymentVersion("GitCommit")
	expected := &DeploymentVersion{
		CommitID: "testing",
	}
	if err != nil {

		t.Fatal(err)

	}

	if out.CommitID != expected.CommitID {

		t.Errorf("output is wrong. Have: %v, want: %v.", out, expected)

	}

}

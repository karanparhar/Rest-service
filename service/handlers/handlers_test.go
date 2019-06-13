package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func TestRouter(t *testing.T) {
	s := mux.NewRouter()
	r := Router(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	os.Setenv("GitCommit", "testing")

	res, err := http.Get(ts.URL + "/api/status")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /api/status is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}

	res, err = http.Post(ts.URL+"/api/status", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code for /api/status is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusMethodNotAllowed)
	}

	res, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /api/status is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusNotFound)
	}
}

func TestDeploymentVersion(t *testing.T) {
	os.Setenv("GitCommit", "testing")
	w := httptest.NewRecorder()
	DeploymentVersion(w, nil)

	resp := w.Result()
	if have, want := resp.StatusCode, http.StatusOK; have != want {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", have, want)
	}
}

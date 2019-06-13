package handlers

import (
	"encoding/json"
	"net/http"

	v "github.com/Rest-service/service/version"
	"github.com/gorilla/mux"
)

// Router register necessary routes and returns an instance of a router.
func Router(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/status", DeploymentVersion).Methods("GET")
	return r
}

func DeploymentVersion(w http.ResponseWriter, r *http.Request) {

	var version = v.NewDeploymentVersion()

	out, err := version.GetDeploymentVersion("GitCommit")

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error: " + err.Error())
		return

	}

	json.NewEncoder(w).Encode(out)

}

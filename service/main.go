package main

import (
	"log"
	"net/http"
	"os"

	h "github.com/Rest-service/service/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	h.Router(r)

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))

}

package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("This is my home page"))
}

func main() {

	// Create the router
	router := mux.NewRouter()

	// Register the route
	home := homeHandler{}

	router.HandleFunc("/", home.ServeHTTP).Methods("GET")

	// Start the server
	http.ListenAndServe(":8010", router)
}

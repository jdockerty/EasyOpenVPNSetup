package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: %d", http.StatusOK)
}

func main() {
	newRouter := mux.NewRouter()

	newRouter.HandleFunc("/", HealthHandler)
	http.ListenAndServe(":8080", newRouter)
}

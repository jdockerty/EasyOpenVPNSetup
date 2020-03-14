package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Client struct {
	ClientName string
	TLSEncrypt string
}

type Status struct {
	Code int
}

// StatusHandler used to check whether the server is running, the 200 OK is returned if it is.
// This provides an easy way to check whether the service is running before constructing a larger POST.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	newStatusResponse := Status{Code: http.StatusOK}
	json.Marshal(newStatusResponse)
	json.NewEncoder(w).Encode(newStatusResponse)
}

// AddClientHandler is used to add a client to the running OpenVPN server.
// It uses the JSON contained within the POST request method to construct a profile on the server. This responds with the constructed profile to insert
// into the OpenVPN Client UI.
func AddClientHandler(w http.ResponseWriter, r *http.Request) {
	var newClient Client
	requestBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(requestBody, &newClient)

	exec.Command("sudo openvpn-install/openvpn-install.sh")
	exec.Command("1")
	exec.Command(newClient.ClientName)
	TLSCommandString := "sudo ls /root/" + newClient.ClientName + ".ovpn"
	
	output, err := exec.Command(TLSCommandString).Output()
	if err != nil {
		panic(err)
	}
	newClient = Client{ClientName : output}

	fmt.Println("New Client:", newClient)
	json.NewEncoder(w).Encode("New Client:" + newClient.ClientName)
}

func main() {
	newRouter := mux.NewRouter()
	fmt.Println("Running...")
	newRouter.HandleFunc("/api/Status", StatusHandler)
	newRouter.HandleFunc("/api/AddClient", AddClientHandler).Methods("POST")
	http.ListenAndServe(":8080", newRouter)
}

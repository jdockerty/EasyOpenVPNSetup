package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os/exec"
	//"bufio"
	//"os"
	"io"
	"bytes"
)

type Client struct {
	ClientName string
	TLSEncrypt []byte
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

	executeOpenVPNScript(newClient.ClientName)

	TLSCommandString := "sudo cat /root/" + newClient.ClientName + ".ovpn"
	
	output, err := exec.Command(TLSCommandString).Output()
	if err != nil {
		panic(err)
	}
	newClient = Client{TLSEncrypt : output}

	fmt.Println("New Client:", newClient)
	json.NewEncoder(w).Encode("New Client:\n" + newClient.ClientName + " " + string(newClient.TLSEncrypt))

}

func executeOpenVPNScript(clientName string) {
	c1 := exec.Command("printf", fmt.Sprintf("1\n%s", clientName))
	c2 := exec.Command("bash", "-c", "sudo ~/openvpn-install/openvpn-install.sh")
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
}

func main() {
	newRouter := mux.NewRouter()
	fmt.Println("Running...")
	newRouter.HandleFunc("/api/Status", StatusHandler)
	newRouter.HandleFunc("/api/AddClient", AddClientHandler).Methods("POST")
	http.ListenAndServe(":8080", newRouter)

}

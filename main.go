package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	//"os"
)

// Client struct for holding corresponding data within the program.
type Client struct {
	Name       string
	TLSEncrypt []byte
}

// Status struct for status code response in JSON.
type Status struct {
	Code int
}

// StatusHandler used to check whether the server is running, the 200 OK is returned if it is.
// This provides an easy way to check whether the service is running before constructing a larger POST.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	newStatusResponse := Status{Code: http.StatusOK}
	json.Marshal(newStatusResponse)
	json.NewEncoder(w).Encode(newStatusResponse)
	log.Println("Status response.")
}

// AddClientHandler is used to add a client to the running OpenVPN server.
// It uses the JSON contained within the POST request method to construct a profile on the server. This responds with the constructed profile to insert
// into the OpenVPN Client UI.
func AddClientHandler(w http.ResponseWriter, r *http.Request) {
	var newClient Client

	requestBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(requestBody, &newClient)

	log.Println("Client received: " + newClient.Name)
	executeOpenVPNScript(newClient, w)
}

func executeReadNewProfile(clientName string) string {
	// Command for reading the .ovpn config file created on the server.
	readConfigCommand := "sudo cat /root/" + clientName + ".ovpn"

	output, err := exec.Command("bash", "-c", readConfigCommand).Output()
	if err != nil {
		panic(err)
	}
	log.Println(".ovpn file read.")
	return string(output)
}

func executeOpenVPNScript(clientToAdd Client, responseWriter http.ResponseWriter) {
	// Command to pipe into the shell script, selects option 1 and adds given client name.
	c1 := exec.Command("printf", fmt.Sprintf("1\n%s", clientToAdd.Name))
	c2 := exec.Command("bash", "-c", "sudo ~/openvpn-install/openvpn-install.sh")

	r, w := io.Pipe()
	c1.Stdout = w // Reader is tied to Stdout of command 1
	c2.Stdin = r  // Writer is tied to Stdin of command 2

	err := c1.Start() // Start command 1 execution
	if err != nil {
		panic(err)
	}

	err = c2.Start() // Execute command 2
	if err != nil {
		panic(err)
	}

	err = c1.Wait()
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	err = c2.Wait()
	if err != nil {
		panic(err)
	}

	clientResponseData := executeReadNewProfile(clientToAdd.Name)
	log.Println("Both commands executed")

	// Write the response to user, allows them to copy/paste the output into an .ovpn file
	responseWriter.Write([]byte(string("Paste the following into an .ovpn file: \n" + clientResponseData)))
	log.Println("Client added and response sent.")
}

// Function used to create a logging file called info.log, log messages are sent here with a timestamp.
// func log.Println((logMessage... string) {
// 	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
//     if err != nil {
//         log.Fatal(err)
//     }

// 	defer logFile.Close()

// 	log.SetOutput(logFile)
// 	log.Println(logMessage)
// }

func main() {
	newRouter := mux.NewRouter()
	log.Println("Running Server...")
	newRouter.HandleFunc("/api/status", StatusHandler)
	newRouter.HandleFunc("/api/addclient", AddClientHandler).Methods("POST")
	http.ListenAndServe(":8080", newRouter)
}

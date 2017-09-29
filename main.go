package main

import (
	"flag"
	"fmt"
	"net/http"
)

// Flag vars
var bindPort int
var accessEmail string
var accessToken string

func init() {
	flag.IntVar(&bindPort, "bind-port", 4567, "the port the HTTP server binds to")
	flag.StringVar(&accessEmail, "access-email", "", "the email account used for auth")
	flag.StringVar(&accessToken, "access-token", "", "the API token used for auth")

	flag.Parse()
}

func main() {
	portString := fmt.Sprintf(":%d", bindPort)
	fmt.Println("[-] Starting on ", portString)
	http.Handle("/receive", receiveHandler())
	http.ListenAndServe(portString, nil)
}

package main

import (
	"flag"
	"fmt"
	"net/http"
)

// Flag vars
var bindPort int

func init() {
	flag.IntVar(&bindPort, "bind-port", 4567, "the port the HTTP server binds to")

	flag.Parse()
}

func main() {
	portString := fmt.Sprintf(":%d", bindPort)
	fmt.Println("[-] Starting on ", portString)
	http.Handle("/receive", receiveHandler())
	http.ListenAndServe(portString, nil)
}

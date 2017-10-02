package main

import (
	"fmt"
	"net/http"

	"github.com/solarwinds/p2l/config"
)

func main() {
	portString := fmt.Sprintf(":%d", config.BindPort())
	fmt.Println("[-] Starting on ", portString)
	http.Handle("/receive", receiveHandler())
	http.ListenAndServe(portString, nil)
}

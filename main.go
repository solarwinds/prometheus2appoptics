package main

import (
	"fmt"
	"net/http"

	"github.com/solarwinds/p2l/config"
	"github.com/solarwinds/p2l/librato"
)

func main() {
	portString := fmt.Sprintf(":%d", config.BindPort())
	fmt.Println("[-] Starting on ", portString)

	lc := librato.NewClient(config.AccessEmail(), config.AccessToken())

	http.Handle("/receive", receiveHandler(lc))
	http.Handle("/spaces", listSpacesHandler(lc))
	http.Handle("/test", testMetricHandler(lc))

	http.ListenAndServe(portString, nil)
}

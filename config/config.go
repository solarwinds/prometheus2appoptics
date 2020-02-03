package config

// This package provides very simple configuration semantics and is designed to be a dependency
// for the rest of the app and any of its component packages.

import (
	"flag"
	"os"
)

// AppName app metadata
const AppName = "prometheus2appoptics"

var (
	globalConf *Config
	bindPort   int
)

func init() {
	flag.IntVar(&bindPort, "bind-port", 4567, "the port the HTTP server binds to")
	flag.Parse()
	globalConf = New()
}

// Config for the reporter
type Config struct {
	bindPort    int
	accessEmail string
}

// New *Config constructor
func New() *Config {
	return &Config{
		bindPort: bindPort,
	}
}

// AccessToken returns a string representing a AppOptics API token
func AccessToken() string {
	return os.Getenv("APPOPTICS_TOKEN")
}

// BindPort returns the port number that the service is bound to
func BindPort() int {
	return globalConf.bindPort
}

// PushErrorLimit is a hardcoded limit on how many errors will be tolerated before the service stops attempting push
func PushErrorLimit() int {
	return 5
}

// SendStats returns true if the application should persist stats over the network to AppOptics, false otherwise
func SendStats() bool {
	return os.Getenv("SEND_STATS") != ""
}

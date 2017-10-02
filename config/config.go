package config

// This package provides very simple configuration semantics and is designed to be a dependency
// for the rest of the app and any of its component packages.

import "flag"

// globalConf is the Config singleton
var globalConf *Config

// Flag vars
var bindPort int
var accessEmail string
var accessToken string
var sendStats bool

func init() {
	flag.IntVar(&bindPort, "bind-port", 4567, "the port the HTTP server binds to")
	flag.StringVar(&accessEmail, "access-email", "", "the email account used for auth")
	flag.StringVar(&accessToken, "access-token", "", "the API token used for auth")
	flag.BoolVar(&sendStats, "send-stats", false, "sends data on the wire if true, prints to stdout if false")

	flag.Parse()

	globalConf = New()
}

type Config struct {
	bindPort    int
	accessEmail string
	accessToken string
	sendStats   bool
}

func New() *Config {
	return &Config{
		bindPort:    bindPort,
		accessEmail: accessEmail,
		accessToken: accessToken,
		sendStats:   sendStats,
	}
}

// AccessEmail returns a string representing an email associated with a Librato API token
func AccessEmail() string {
	return globalConf.accessEmail
}

// AccessToken returns a string representing a Librato API token
func AccessToken() string {
	return globalConf.accessToken
}

// BindPort returns the port number that the service is bound to
func BindPort() int {
	return globalConf.bindPort
}

// PushErrorLimit is a hardcoded limit on how many errors will be tolerated before the service stops attempting push
func PushErrorLimit() int {
	return 5
}

// SendStats returns true if the application should persist stats over the network to Librato, false otherwise
func SendStats() bool {
	return globalConf.sendStats
}

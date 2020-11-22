package agent

import (
	"log"
	"time"
)

// Config defines the configurable options for an Agent. Ultimately, binaries
// that implement an Agent will be dynamically generated using the text/template
// package. The Config structure will be passed to the template to specify the
// initial values for that current Agent.
type Config struct {
	// server information
	Host string
	Port int

	// agent configurationinformation
	FetchFrequency         string
	fetchFrequencyDuration time.Duration
}

func (a *Agent) host() string {
	return a.config.Host
}

func (a *Agent) setHost(host string) {
	log.Printf("setting host to %s\n", host)
	a.config.Host = host
}

func (a *Agent) port() int {
	return a.config.Port
}

func (a *Agent) setPort(port int) {
	log.Printf("setting port to %d\n", port)
	a.config.Port = port
}

func (a *Agent) fetchFrequency() string {
	return a.config.FetchFrequency
}

// setFetchFrequency modifies the duration of time between calls to FetchCommand().
// The input to this function is a string indicating a duration (i.e. "1m", "5s",
// etc). If there is an error parsing the provided duration, that error is
// returned and the value of FetchFrequency is not updated.
func (a *Agent) setFetchFrequency(duration string) error {
	// if the provided duration string is an invalid duration, nothing should
	// be updated.
	d, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	log.Printf("setting fetchFrequency to %s\n", duration)
	a.config.FetchFrequency = duration
	a.config.fetchFrequencyDuration = d
	return nil
}

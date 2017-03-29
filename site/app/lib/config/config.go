package config

import (
	"io"
	"io/ioutil"
	"os"
)

// Parser must implement ParseJSON
type Parser interface {
	ParseJSON([]byte) error
}

// Load the JSON config file
func Load(configFile string, p Parser) error {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		return err
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		return err
	}

	// Parse the config
	if err := p.ParseJSON(jsonBytes); err != nil {
		return err
	}

	return nil
}

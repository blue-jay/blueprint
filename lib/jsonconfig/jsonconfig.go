package jsonconfig

import (
	"io/ioutil"
	"log"
)

// Parser must implement ParseJSON
type Parser interface {
	ParseJSON([]byte) error
}

// LoadOrFatal loads the JSON config file and exits if it can't be parsed
func LoadOrFatal(configFile string, p Parser) {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	err = p.ParseJSON(jsonBytes)
	if err != nil {
		log.Fatalln("Could not parse %q: %v", configFile, err)
	}
}

// Load the JSON config file
func Load(configFile string, p Parser) error {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Parse the config
	if err := p.ParseJSON(jsonBytes); err != nil {
		return err
	}

	return nil
}

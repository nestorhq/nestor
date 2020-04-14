package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ResourcesDefinition list the optional resources that we want in the application
type ResourcesDefinition struct {
	ID string
}

// ApplicationDefinition application definition
type ApplicationDefinition struct {
	Name string
}

// Config nestor configuration
type Config struct {
	Nestor      string
	Application ApplicationDefinition `json:"application"`
	Resources   []ResourcesDefinition `json:"resources"`
}

// ReadConfig read congiration from file
func ReadConfig(filename string) (*Config, error) {
	var config Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("source: %s", source)
	j2, err := yaml.YAMLToJSON(source)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(j2))
	dec.DisallowUnknownFields()

	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

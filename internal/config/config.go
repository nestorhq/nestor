package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// LambdaPermissionAction describe an action allowed
type LambdaPermissionAction struct {
	Operation string
}

// LambdaPermission describe a lambda permission
type LambdaPermission struct {
	ResourceID string                   `json:"resourceId"`
	Actions    []LambdaPermissionAction `json:"actions"`
}

// LambdasDefinition list the optional resources that we want in the application
type LambdasDefinition struct {
	ID          string
	Permissions []LambdaPermission
}

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
	Lambdas     []LambdasDefinition   `json:"lambdas"`
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

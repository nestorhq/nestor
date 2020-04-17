package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// TrigerS3UploadDefinition trigger associated to s3 upload
type TriggerS3UploadDefinition struct {
	Lambda string
	Prefix string
	Suffix string
}

// TriggersDefinition triggers description
type TriggersDefinition struct {
	S3upload []TriggerS3UploadDefinition `json:"s3upload"`
}

// LambdaPermissionAction describe an action allowed
type LambdaPermissionAction struct {
	Operation string
}

// LambdaPermission describe a lambda permission
type LambdaPermission struct {
	ResourceID string                   `json:"resourceId"`
	Actions    []LambdaPermissionAction `json:"actions"`
}

// LambdaDefinition list the optional resources that we want in the application
type LambdaDefinition struct {
	ID          string
	Permissions []LambdaPermission
}

// ResourceDefinition list the optional resources that we want in the application
type ResourceDefinition struct {
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
	Resources   []ResourceDefinition  `json:"resources"`
	Lambdas     []LambdaDefinition    `json:"lambdas"`
	Triggers    []TriggersDefinition  `json:"triggers"`
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

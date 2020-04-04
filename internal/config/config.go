package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ResourceS3Bucket s3 bucket description
type ResourceS3Bucket struct {
	ID string
}

// ResourceDynamoDbTable dynamodb table description
type ResourceDynamoDbTable struct {
	ID string
}

// ResourceHTTPAPI http api description
type ResourceHTTPAPI struct {
	ID string
}

// ResourcesDefinition resource definition description
type ResourcesDefinition struct {
	S3Bucket      []ResourceS3Bucket      `json:"s3Bucket"`
	DynamoDbTable []ResourceDynamoDbTable `json:"dynamoDbTable"`
	HTTPAPI       []ResourceHTTPAPI       `json:"httpApi"`
}

// AppDefinition application definition
type AppDefinition struct {
	Name string
}

// Config nestor configuration
type Config struct {
	Nestor    string
	App       AppDefinition       `json:"app"`
	Resources ResourcesDefinition `json:"resources"`
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

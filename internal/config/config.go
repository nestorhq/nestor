package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type ResourceS3Bucket struct {
	Id string
}

type ResourceDynamoDbTable struct {
	Id string
}

type ResourcesDefinition struct {
	S3Bucket      []ResourceS3Bucket      `json:"s3Bucket"`
	DynamoDbTable []ResourceDynamoDbTable `json:"dynamoDbTable"`
}

type Config struct {
	Version   string
	Resources ResourcesDefinition `json:"resources"`
}

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

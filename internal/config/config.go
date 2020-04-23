package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ResourcesLambdaFunctionDeployment deployment definition
type ResourcesLambdaFunctionDeployment struct {
	ID      string `json:"id"`
	File    string `json:"file"`
	Handler string `json:"handler"`
}

// DeploymentsDefinition deployments descriptions
type DeploymentsDefinition struct {
	LambdaFunction []ResourcesLambdaFunctionDeployment `json:"lambda_function"`
}

// TriggerS3CopyNotificationDefinition trigger associated to s3 upload
type TriggerS3CopyNotificationDefinition struct {
	LambdaID string `json:"lambdaId"`
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
}

// TriggerS3CopyDefinition trigger associated to s3 upload
type TriggerS3CopyDefinition struct {
	BucketID string                                `json:"bucketId"`
	Lambdas  []TriggerS3CopyNotificationDefinition `json:"lambdas"`
}

// TriggersDefinition triggers description
type TriggersDefinition struct {
	S3copy []TriggerS3CopyDefinition `json:"s3copy"`
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

// ResourcesLambdaFunctionDefinition list the optional resources that we want in the application
type ResourcesLambdaFunctionDefinition struct {
	ID          string
	Runtime     string
	Permissions []LambdaPermission
}

// ResourceCognitoUserPoolDefinition resource definition
type ResourceCognitoUserPoolDefinition struct {
	ID string
}

// ResourceCognitoDefinition definition of resource
type ResourceCognitoDefinition struct {
	UserPool ResourceCognitoUserPoolDefinition
}

// ResourceDynamodbTableDefinition resource definition
type ResourceDynamodbTableDefinition struct {
	ID string
}

// ResourceDynamodbDefinition definition of resource
type ResourceDynamodbDefinition struct {
	Table ResourceDynamodbTableDefinition
}

// ResourceS3BucketDefinition resource definition
type ResourceS3BucketDefinition struct {
	ID         string
	BucketName string
}

// ResourcesCloudwatchLogsGroupDefinition resource definition
type ResourcesCloudwatchLogsGroupDefinition struct {
	ID string
}

// ResourcesEventBridgeBusDefinition resource definition
type ResourcesEventBridgeBusDefinition struct {
	ID string
}

// ResourcesApigatewayHTTPDefinition resource definition
type ResourcesApigatewayHTTPDefinition struct {
	ID             string
	TargetLambdaID string
}

// ResourceDefinition list the optional resources that we want in the application
type ResourceDefinition struct {
	CognitoUserPool     []ResourceCognitoUserPoolDefinition      `json:"cognito_userpool"`
	DynamodbTable       []ResourceDynamodbTableDefinition        `json:"dynamodb_table"`
	S3Bucket            []ResourceS3BucketDefinition             `json:"s3_bucket"`
	LambdaFunction      []ResourcesLambdaFunctionDefinition      `json:"lambda_function"`
	CloudwatchlogsGroup []ResourcesCloudwatchLogsGroupDefinition `json:"cloudwatchlogs_group"`
	EventBridgeBus      []ResourcesEventBridgeBusDefinition      `json:"eventbridge_bus"`
	ApigatewayHTTP      []ResourcesApigatewayHTTPDefinition      `json:"apigateway_http"`
}

// ApplicationDefinition application definition
type ApplicationDefinition struct {
	Name          string
	ProfileName   string // name of profile to use to initialize aws API
	Region        string // default region for resources creation
	RegionCognito string // region for cognito resource
}

// Config nestor configuration
type Config struct {
	Nestor      string
	Application ApplicationDefinition `json:"application"`
	Resources   ResourceDefinition    `json:"resources"`
	Triggers    TriggersDefinition    `json:"triggers"`
	Deployments DeploymentsDefinition `json:"deployments"`
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

package config

import "errors"

// ResCognitoMain resource id
const ResCognitoMain = "nestor.res.cognito.main"

// ResDynamoDbTableMain resource id
const ResDynamoDbTableMain = "nestor.res.dynamoDbTable.main"

// ResEventBridgeMain resource id
const ResEventBridgeMain = "nestor.res.eventBridge.main"

// ResS3BucketForStorage resource id
const ResS3BucketForStorage = "nestor.res.s3.storage"

// ResS3BucketForUpload reource id
const ResS3BucketForUpload = "nestor.res.s3.upload"

// ResHTTPAPIMain resource id
const ResHTTPAPIMain = "nestor.res.httpApi.main"

// ResLogGroupMainEventBridge resource id
const ResLogGroupMainEventBridge = "nestor.res.logGroup.mainEventBridgeTarget"

// ResourceType defines the type of a resource
type ResourceType int

const (
	unknown ResourceType = iota
	cognitoUserPool
	dynamoDbTable
	eventBridgeBus
	s3Bucket
	httpAPIGateway
	cloudwatchLogGroup
	lambda
)

// ResourceDescription resource description
type ResourceDescription struct {
	ID           string // the id of the resource
	description  string
	isOptional   bool
	resourceType ResourceType
}

type registeredResource struct {
	resourceID   string
	awsID        string
	resourceType ResourceType
}

// structure to store the registered resources
var _registeredResources = map[string]registeredResource{}

var _nestorResources = []ResourceDescription{
	{
		ID:           ResCognitoMain,
		description:  "the application cognito user pool",
		isOptional:   false,
		resourceType: cognitoUserPool,
	},
	{
		ID:           ResDynamoDbTableMain,
		description:  "the main dynamoDb table (following a single table pattern)",
		isOptional:   false,
		resourceType: dynamoDbTable,
	},
	{
		ID:           ResEventBridgeMain,
		description:  "the main event bridge bus",
		isOptional:   false,
		resourceType: eventBridgeBus,
	},
	{
		ID:           ResS3BucketForStorage,
		description:  "the s3 bucket name for storage",
		isOptional:   false,
		resourceType: s3Bucket,
	},
	{
		ID:           ResS3BucketForUpload,
		description:  "the s3 bucket name for upload",
		isOptional:   false,
		resourceType: s3Bucket,
	},
	{
		ID:           ResHTTPAPIMain,
		description:  "the Api Gateway main http service",
		isOptional:   false,
		resourceType: httpAPIGateway,
	},
	{
		ID:           ResLogGroupMainEventBridge,
		description:  "the Cloudwatch group name to push event from main EventBridge",
		isOptional:   true,
		resourceType: cloudwatchLogGroup,
	},
}

// IsResourceRequired indicates if a given resource must be provisioned
func (config *Config) IsResourceRequired(resourceID string) (bool, string) {
	// check first if that's a nestor required resource
	for _, nestorResource := range _nestorResources {
		if nestorResource.ID == resourceID {
			if !nestorResource.isOptional {
				return true, resourceID
			}
		}
	}

	// check if the user requested that resource
	for _, resource := range config.Resources {
		if resource.ID == resourceID {
			return true, resourceID
		}
	}
	return false, resourceID
}

func findNestorResourceByID(resourceID string) *ResourceDescription {
	for _, nestorResource := range _nestorResources {
		if nestorResource.ID == resourceID {
			return &nestorResource
		}
	}
	return nil
}

// RegisterNestorResource register a resource with its arn
func (config *Config) RegisterNestorResource(resourceID string, awsID string) error {
	resource := findNestorResourceByID(resourceID)
	if resource == nil {
		return errors.New("unknown resource:" + resourceID)
	}
	_registeredResources[resourceID] = registeredResource{
		awsID:        awsID,
		resourceID:   resourceID,
		resourceType: resource.resourceType,
	}
	return nil
}

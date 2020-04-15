package config

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

// ResourceDescription resource description
type ResourceDescription struct {
	ID          string // the id of the resource
	description string
	isOptional  bool
}

var availableResources = []ResourceDescription{
	{
		ID:          ResCognitoMain,
		description: "the application cognito user pool",
		isOptional:  false,
	},
	{
		ID:          ResDynamoDbTableMain,
		description: "the main dynamoDb table (following a single table pattern)",
		isOptional:  false,
	},
	{
		ID:          ResEventBridgeMain,
		description: "the main event bridge bus",
		isOptional:  false,
	},
	{
		ID:          ResS3BucketForStorage,
		description: "the s3 bucket name for storage",
		isOptional:  false,
	},
	{
		ID:          ResS3BucketForUpload,
		description: "the s3 bucket name for upload",
		isOptional:  false,
	},
	{
		ID:          ResHTTPAPIMain,
		description: "the Api Gateway main http service",
		isOptional:  false,
	},
	{
		ID:          ResLogGroupMainEventBridge,
		description: "the Cloudwatch group name to push event from main EventBridge",
		isOptional:  true,
	},
}

// IsResourceRequired indicates if a given resource must be provisioned
func (config *Config) IsResourceRequired(resourceID string) (bool, string) {
	// check first if that's a nestor required resource
	for _, availableResource := range availableResources {
		if availableResource.ID == resourceID {
			if !availableResource.isOptional {
				return false, resourceID
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

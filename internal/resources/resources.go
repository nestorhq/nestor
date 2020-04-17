package resources

import (
	"github.com/nestorhq/nestor/internal/config"
)

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

// ResourceAttName type for resource attributes names
type ResourceAttName string

const (
	// AttArn arn attribute
	AttArn ResourceAttName = "arn"
	// AttID id attribute
	AttID ResourceAttName = "id"
	// AttName name attribute
	AttName ResourceAttName = "name"
)

// ResourceDescription resource description
type ResourceDescription struct {
	ID           string // the id of the resource
	description  string
	isOptional   bool
	resourceType ResourceType
}

// RegisteredResource decsription of a registered resources
type RegisteredResource struct {
	resourceID   string
	attributes   map[ResourceAttName]string
	resourceType ResourceType
}

// GetAttribute retrieve attribute
func (res *RegisteredResource) GetAttribute(name ResourceAttName) string {
	return res.attributes[name]
}

// Resources hold the resources
type Resources struct {
	// structure to store the registered resources
	registeredResources map[string]*RegisteredResource
	nestorResources     []ResourceDescription
}

// NewResources ctor
func NewResources() *Resources {
	var nestorResources = []ResourceDescription{
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
	var result = Resources{
		registeredResources: make(map[string]*RegisteredResource),
		nestorResources:     nestorResources,
	}

	return &result
}

// IsResourceRequired indicates if a given resource must be provisioned
func (res *Resources) IsResourceRequired(resourceID string, resDef []config.ResourceDefinition) (bool, string) {
	// check first if that's a nestor required resource
	for _, nestorResource := range res.nestorResources {
		if nestorResource.ID == resourceID {
			if !nestorResource.isOptional {
				return true, resourceID
			}
		}
	}

	// check if the user requested that resource
	for _, resource := range resDef {
		if resource.ID == resourceID {
			return true, resourceID
		}
	}
	return false, resourceID
}

func (res *Resources) findNestorResourceByID(resourceID string) *ResourceDescription {
	for _, nestorResource := range res.nestorResources {
		if nestorResource.ID == resourceID {
			return &nestorResource
		}
	}
	return nil
}

// RegisterNestorResource register a resource with its arn
func (res *Resources) RegisterNestorResource(resourceID string, attName ResourceAttName, attValue string) error {
	registered, ok := res.registeredResources[resourceID]
	if !ok {
		registered = &RegisteredResource{
			resourceID:   resourceID,
			resourceType: unknown, // TODO
		}
		registered.attributes = make(map[ResourceAttName]string)
		res.registeredResources[resourceID] = registered
	}
	var attributes = registered.attributes

	attributes[attName] = attValue

	return nil
}

// FindResourceByID find creation information about a given resource
func (res *Resources) FindResourceByID(resourceID string) *RegisteredResource {
	resource, ok := res.registeredResources[resourceID]
	if ok {
		return resource
	}
	return nil
}

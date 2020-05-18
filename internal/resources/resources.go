package resources

import (
	"github.com/nestorhq/nestor/internal/config"
)

// ResourceType defines the type of a resource
type ResourceType int

const (
	unknown ResourceType = iota
	// CognitoUserPool type
	CognitoUserPool
	// DynamoDbTable type
	DynamoDbTable
	//EventBridgeBus type
	EventBridgeBus
	//S3Bucket type
	S3Bucket
	// HTTPAPIGateway type
	HTTPAPIGateway
	// CloudwatchLogGroup type
	CloudwatchLogGroup
	// LambdaFunction type
	LambdaFunction
	// SESMail send mail
	SESMail
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
	var result = Resources{
		registeredResources: make(map[string]*RegisteredResource),
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
func (res *Resources) RegisterNestorResource(resourceID string, resourceType ResourceType, attName ResourceAttName, attValue string) error {
	registered, ok := res.registeredResources[resourceID]
	if !ok {
		registered = &RegisteredResource{
			resourceID:   resourceID,
			resourceType: resourceType,
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

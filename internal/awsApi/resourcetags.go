package awsapi

// ResourceTags tags associated to created resources
type ResourceTags struct {
	environment string
	appName     string
}

// NewResourceTag constructor
func NewResourceTag(environment string, appName string) *ResourceTags {
	var tag = ResourceTags{}
	tag.appName = appName
	tag.environment = environment
	return &tag
}

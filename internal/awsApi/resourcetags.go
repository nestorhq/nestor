package awsapi

import "errors"

// ResourceTags tags associated to created resources
type ResourceTags struct {
	nestorVersion string
	environment   string
	appName       string
}

// NewResourceTag constructor
func NewResourceTag(nestorVersion string, environment string, appName string) *ResourceTags {
	var tag = ResourceTags{}
	tag.appName = appName
	tag.environment = environment
	tag.nestorVersion = nestorVersion
	return &tag
}

func (t *ResourceTags) getTagsAsMap() map[string]string {
	var tags = map[string]string{
		"appName":     t.appName,
		"environment": t.environment,
		"nv":          t.nestorVersion,
	}
	return tags
}

func (t *ResourceTags) checkTags(tags map[string]*string) error {
	if *tags["appName"] != t.appName {
		return errors.New("resource exist with bad tag(appName)")
	}
	if *tags["environment"] != t.environment {
		return errors.New("resource exist with bad tag(environment)")
	}
	return nil
}

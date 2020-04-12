package awsapi

import (
	"errors"
)

// ResourceTags tags associated to created resources
type ResourceTags struct {
	nestorVersion string
	environment   string
	appName       string
}

// Tag tag description
type Tag struct {
	Key   string
	Value string
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

func (t *ResourceTags) getTagsAsMapWithID(id string) map[string]string {
	var tags = t.getTagsAsMap()
	tags["id"] = id
	return tags
}

func (t *ResourceTags) getTagsAsTags() []Tag {
	var result = make([]Tag, 0, 4)
	result = append(result, Tag{
		Key:   "appName",
		Value: t.appName,
	}, Tag{
		Key:   "environment",
		Value: t.environment,
	}, Tag{
		Key:   "nv",
		Value: t.nestorVersion,
	})
	return result
}

func (t *ResourceTags) getTagsAsTagsWithID(id string) []Tag {
	result := t.getTagsAsTags()
	result = append(result, Tag{
		Key:   "id",
		Value: id,
	})
	return result
}

func (t *ResourceTags) checkTags(tags map[string]*string, id string) error {
	if *tags["appName"] != t.appName {
		return errors.New("resource exist with bad tag(appName)")
	}
	if *tags["environment"] != t.environment {
		return errors.New("resource exist with bad tag(environment)")
	}
	if *tags["id"] != id {
		return errors.New("resource exist with bad tag(environment)")
	}
	return nil
}

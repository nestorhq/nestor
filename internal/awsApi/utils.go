package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
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
	tags["nestorId"] = id
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
		Key:   "nestorId",
		Value: id,
	})
	return result
}

func (t *ResourceTags) checkTagValue(tags map[string]*string, tagName string, expected string) error {
	if pval, ok := tags[tagName]; ok {
		if *pval != expected {
			return fmt.Errorf("resource exist with bad tag(%s) expected: %s, actual: %s", tagName, expected, *pval)
		}
	} else {
		return fmt.Errorf("missing tag (%s) expected: %s", tagName, expected)
	}
	return nil
}

func (t *ResourceTags) checkTags(tags map[string]*string, id string) error {
	var err error
	err = t.checkTagValue(tags, "appName", t.appName)
	if err != nil {
		return err
	}
	err = t.checkTagValue(tags, "environment", t.environment)
	if err != nil {
		return err
	}
	err = t.checkTagValue(tags, "nestorId", id)
	if err != nil {
		return err
	}
	return nil
}

func getAwsErrorCode(err error) string {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code()
	}
	return ""
}

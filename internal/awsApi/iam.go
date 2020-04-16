package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/nestorhq/nestor/internal/reporter"
)

// IAMAPI api
type IAMAPI struct {
	resourceTags *ResourceTags
	client       *iam.IAM
}

// IAMInformation description of a IAM
type IAMInformation struct {
	functionName string
	arn          string
}

// RoleInformation  description of an IAM role
type RoleInformation struct {
	RoleArn  string
	RoleName string
}

// NewIAMAPI constructor
func NewIAMAPI(session *session.Session, resourceTags *ResourceTags) (*IAMAPI, error) {
	var api = IAMAPI{resourceTags: resourceTags}
	// Create IAM client
	api.client = iam.New(session)
	return &api, nil
}

// this is used to allow the function to be called
// by the lambda service
const assumePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

func (api *IAMAPI) checkRoleExistenceAndTags(roleName string, nestorID string, t *reporter.Task) (*RoleInformation, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.GetRole").WithArg("roleName", roleName))

	tags := api.resourceTags.getTagsAsTagsWithID(nestorID)
	iamTags := make([]*iam.Tag, 0, 4)
	for _, t := range tags {
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}

	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}
	result, err := api.client.GetRole(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			return nil, nil
		}
		t0.Fail(err)
		return nil, err
	}
	t0.LogM(reporter.NewMessage("GetRole result").
		WithArg("input", input.GoString()).
		WithArg("result", result.GoString()))

	return &RoleInformation{
		RoleArn:  *result.Role.Arn,
		RoleName: *result.Role.RoleName,
	}, nil
}

func (api *IAMAPI) doCreateRole(roleName string, nestorID string, t *reporter.Task) (*RoleInformation, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.CreateRole").WithArg("roleName", roleName))

	tags := api.resourceTags.getTagsAsTagsWithID(nestorID)
	iamTags := make([]*iam.Tag, 0, 4)
	for _, t := range tags {
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumePolicy),
		Description:              aws.String("Role for lambda"),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags:                     iamTags,
	}
	result, err := api.client.CreateRole(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	t0.LogM(reporter.NewMessage("CreateRole result").
		WithArg("input", input.GoString()).
		WithArg("result", result.GoString()))

	// check tags
	tagsToCheck := map[string]*string{}
	tagsFromRole := result.Role.Tags
	for _, tag := range tagsFromRole {
		tagsToCheck[*tag.Key] = tag.Value
	}

	t1 := t0.SubM(reporter.NewMessage("checking tags").WithArgs(tagsToCheck))
	err2 := api.resourceTags.checkTags(tagsToCheck, nestorID)
	if err2 != nil {
		t1.Fail(err2)
		return nil, err2
	}

	return &RoleInformation{
		RoleArn:  *result.Role.Arn,
		RoleName: *result.Role.RoleName,
	}, nil
}

// CreateRole create role
func (api *IAMAPI) CreateRole(roleName string, nestorID string, t *reporter.Task) (*RoleInformation, error) {
	t0 := t.SubM(reporter.NewMessage("CreateRole").WithArg("roleName", roleName))

	t1 := t0.Sub("check if role exists")

	result, err := api.checkRoleExistenceAndTags(roleName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}
	if result != nil {
		t1.Log("Role exists")
		t1.Ok()
		return result, nil
	}
	t1.Log("Role does not exist - cretaing it")
	result, err = api.doCreateRole(roleName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		t0.Fail(err)
		return nil, err
	}
	return result, err
}

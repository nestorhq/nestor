package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/nestorhq/nestor/internal/reporter"
)

// CloudWatchLogsAPI api
type CloudWatchLogsAPI struct {
	resourceTags *ResourceTags
	client       *cloudwatchlogs.CloudWatchLogs
}

// CloudWatchLogGroupInformation description of a CloudWatch
type CloudWatchLogGroupInformation struct {
	GroupName string
}

// NewCloudWatchLogsAPI constructor
func NewCloudWatchLogsAPI(session *session.Session, resourceTags *ResourceTags) (*CloudWatchLogsAPI, error) {
	var api = CloudWatchLogsAPI{resourceTags: resourceTags}
	// Create CloudWatch client
	api.client = cloudwatchlogs.New(session)
	return &api, nil
}

func (api *CloudWatchLogsAPI) checkLogGroupTags(groupName string, nestorID string, task *reporter.Task) (bool, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.ListTagsLogGroup").WithArg("groupName", groupName))
	input := &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: aws.String(groupName),
	}
	result, err := api.client.ListTagsLogGroup(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			return false, nil
		}
		t0.Fail(err)
		return false, err
	}
	// t0.LogM(reporter.NewMessage("tags for resource").
	// 	WithArg("input", input.GoString()).
	// 	WithArg("result", result.GoString()))

	// check tags
	t1 := task.SubM(reporter.NewMessage("checkTags").WithArgs(result.Tags))
	err2 := api.resourceTags.checkTags(result.Tags, nestorID)
	if err2 != nil {
		t1.Fail(err2)
		return false, err2
	}
	t1.Ok()
	return true, nil
}

func (api *CloudWatchLogsAPI) tagLogGroup(groupName string, nestorID string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("api.client.TagLogGroup").WithArg("groupName", groupName))
	input := &cloudwatchlogs.TagLogGroupInput{
		LogGroupName: aws.String(groupName),
		Tags:         aws.StringMap(api.resourceTags.getTagsAsMapWithID(nestorID)),
	}
	result, err := api.client.TagLogGroup(input)
	if err != nil {
		t0.Fail(err)
		return err
	}
	t0.LogM(reporter.NewMessage("TagLogGroup operation").
		WithArg("input", input.GoString()).
		WithArg("result", result.GoString()))

	t0.Ok()
	return nil
}

func (api *CloudWatchLogsAPI) doCreateLogGroup(groupName string, nestorID string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("api.client.CreateLogGroup").WithArg("groupName", groupName))
	input := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(groupName),
	}
	_, err := api.client.CreateLogGroup(input)
	if err != nil {
		t0.Fail(err)
		return err
	}
	// t0.LogM(reporter.NewMessage("CreateLogGroup result").
	// 	WithArg("input", input.GoString()).
	// 	WithArg("result", result.GoString()))
	return nil
}

func (api *CloudWatchLogsAPI) createLogGroup(groupName string, nestorID string, t *reporter.Task) (*CloudWatchLogGroupInformation, error) {
	t0 := t.SubM(reporter.NewMessage("createLogGroup").WithArg("groupName", groupName))

	isPresent, err := api.checkLogGroupTags(groupName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if isPresent {
		t0.Log("log group exists")
		t0.Ok()
		return &CloudWatchLogGroupInformation{
			GroupName: groupName,
		}, nil
	}
	t0.Log("log group does not exist")
	err = api.doCreateLogGroup(groupName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}

	err = api.tagLogGroup(groupName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}

	t0.Ok()
	return &CloudWatchLogGroupInformation{
		GroupName: groupName,
	}, nil
}

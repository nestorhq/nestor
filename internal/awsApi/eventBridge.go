package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/nestorhq/nestor/internal/reporter"
)

// EventBridgeAPI api
type EventBridgeAPI struct {
	resourceTags *ResourceTags
	client       *eventbridge.EventBridge
}

// EventBusInformation description of a EventBridge
type EventBusInformation struct {
	eventBusName string
	arn          string
}

// NewEventBridgeAPI constructor
func NewEventBridgeAPI(session *session.Session, resourceTags *ResourceTags) (*EventBridgeAPI, error) {
	var api = EventBridgeAPI{resourceTags: resourceTags}
	// Create EventBridge client
	api.client = eventbridge.New(session)
	return &api, nil
}

func (api *EventBridgeAPI) doCreateEventBus(eventBusName string, nestorID string, task *reporter.Task) (*EventBusInformation, error) {
	t0 := task.SubM(reporter.NewMessage("eventbridge.CreateEventBus").WithArg("eventBusName", eventBusName))

	tags := api.resourceTags.getTagsAsTagsWithID(nestorID)
	eventBusTags := make([]*eventbridge.Tag, 0, 4)
	for _, t := range tags {
		eventBusTags = append(eventBusTags, &eventbridge.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}

	input := &eventbridge.CreateEventBusInput{
		Name: &eventBusName,
		Tags: eventBusTags,
	}

	result, err := api.client.CreateEventBus(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	// fmt.Printf("result: %v\n", result)
	return &EventBusInformation{
		arn:          *result.EventBusArn,
		eventBusName: eventBusName,
	}, nil
}

func (api *EventBridgeAPI) checkEventBusExistence(eventBusName string, task *reporter.Task) (*EventBusInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.DescribeEventBus").WithArg("eventBusName", eventBusName))
	input := &eventbridge.DescribeEventBusInput{
		Name: aws.String(eventBusName),
	}
	result, err := api.client.DescribeEventBus(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			return nil, nil
		}
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"Name": *result.Name,
		"Arn":  *result.Arn,
	})

	return &EventBusInformation{
		eventBusName: *result.Name,
		arn:          *result.Arn,
	}, nil
}

func (api *EventBridgeAPI) checkEventBusTags(eventBusArn string, id string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("api.client.ListTagsOfResource").WithArg("eventBusArn", eventBusArn))
	input := &eventbridge.ListTagsForResourceInput{
		ResourceARN: aws.String(eventBusArn),
	}
	result, err := api.client.ListTagsForResource(input)
	if err != nil {
		t0.Fail(err)
		return err
	}

	tagsToCheck := map[string]*string{}
	tags := result.Tags
	for _, tag := range tags {
		tagsToCheck[*tag.Key] = tag.Value
	}
	// check tags
	t1 := task.SubM(reporter.NewMessage("checkTags").WithArgs(tagsToCheck))
	err2 := api.resourceTags.checkTags(tagsToCheck, id)
	if err2 != nil {
		t1.Fail(err2)
		return err2
	}
	t1.Ok()
	return nil
}

func (api *EventBridgeAPI) checkEventBusExistenceAndTags(eventBusName string, id string, task *reporter.Task) (*EventBusInformation, error) {
	t0 := task.SubM(reporter.NewMessage("checkEventBusExistenceAndTags").WithArg("eventBusName", eventBusName))
	eventBusInformation, err := api.checkEventBusExistence(eventBusName, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if eventBusInformation == nil {
		t0.Ok()
		return nil, nil
	}

	t1 := task.SubM(reporter.NewMessage("checkEventBusTags").WithArg("eventBusName", eventBusName))
	err2 := api.checkEventBusTags(eventBusInformation.arn, id, t1)
	if err2 != nil {
		t1.Fail(err2)
		return nil, err2
	}
	return eventBusInformation, nil
}

func (api *EventBridgeAPI) createEventBus(eventBusName string, id string, task *reporter.Task) (*EventBusInformation, error) {
	t0 := task.SubM(reporter.NewMessage("createEventBus").WithArg("eventBusName", eventBusName))

	t1 := t0.Sub("check if event bus exists")
	eventBusInformation, err := api.checkEventBusExistenceAndTags(eventBusName, id, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}

	if eventBusInformation != nil {
		t1.Log("event bus exists")
		t1.Okr(map[string]string{
			"arn":          eventBusInformation.arn,
			"eventBusName": eventBusInformation.eventBusName,
		})

		return eventBusInformation, nil
	}

	t2 := t0.Sub("event bus does not exist - creating it")
	result, err := api.doCreateEventBus(eventBusName, id, t2)
	if err != nil {
		t2.Fail(err)
	}
	t2.Ok()
	t0.Okr(map[string]string{
		"arn":          result.arn,
		"eventBusName": result.eventBusName,
	})
	return result, nil
}

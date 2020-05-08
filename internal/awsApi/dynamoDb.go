package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nestorhq/nestor/internal/reporter"
)

// TableInformation description of a user pool
type TableInformation struct {
	TableID   string
	TableArn  string
	TableName string
}

// DynamoDbAPI Access to DynamoDb API
type DynamoDbAPI struct {
	resourceTags *ResourceTags
	client       *dynamodb.DynamoDB
}

// NewDynamoDbAPI constructor
func NewDynamoDbAPI(session *session.Session, resourceTags *ResourceTags) (*DynamoDbAPI, error) {
	var api = DynamoDbAPI{resourceTags: resourceTags}
	// Create DynamoDB client
	api.client = dynamodb.New(session)
	return &api, nil
}

func (api *DynamoDbAPI) doCreateMonoTable(tableName string, nestorID string, task *reporter.Task) (*TableInformation, error) {
	t0 := task.SubM(reporter.NewMessage("dynamodb.CreateTableInput").WithArg("tableName", tableName))

	tags := api.resourceTags.getTagsAsTagsWithID(nestorID)
	dynamodbTags := make([]*dynamodb.Tag, 0, 4)
	for _, t := range tags {
		dynamodbTags = append(dynamodbTags, &dynamodb.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("gsi_1"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("sk"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("pk"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(tableName),
		Tags:        dynamodbTags,
	}

	result, err := api.client.CreateTable(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	// fmt.Printf("result: %v\n", result)
	return &TableInformation{
		TableID:   *result.TableDescription.TableId,
		TableArn:  *result.TableDescription.TableArn,
		TableName: *result.TableDescription.TableName,
	}, nil
}

func (api *DynamoDbAPI) checkTableTags(tableArn string, nestorID string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("api.client.ListTagsOfResource").WithArg("tableArn", tableArn))
	input := &dynamodb.ListTagsOfResourceInput{
		ResourceArn: aws.String(tableArn),
	}
	result, err := api.client.ListTagsOfResource(input)
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
	err2 := api.resourceTags.checkTags(tagsToCheck, nestorID)
	if err2 != nil {
		t1.Fail(err2)
		return err2
	}
	t1.Ok()
	return nil
}

func (api *DynamoDbAPI) checkTableExistence(tableName string, task *reporter.Task) (*TableInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.DescribeTable").WithArg("tableName", tableName))
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	result, err := api.client.DescribeTable(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			return nil, nil
		}
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"TableId":   *result.Table.TableId,
		"TableArn":  *result.Table.TableArn,
		"TableName": *result.Table.TableName,
	})

	return &TableInformation{
		TableID:   *result.Table.TableId,
		TableArn:  *result.Table.TableArn,
		TableName: *result.Table.TableName,
	}, nil
}

func (api *DynamoDbAPI) checkTableExistenceAndTags(tableName string, nestorID string, task *reporter.Task) (*TableInformation, error) {
	t0 := task.SubM(reporter.NewMessage("checkTableExistenceAndTags").WithArg("tableName", tableName))
	tableInformation, err := api.checkTableExistence(tableName, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if tableInformation == nil {
		t0.Ok()
		return nil, nil
	}

	t1 := task.SubM(reporter.NewMessage("checkTableTags").WithArg("tableName", tableName))
	err2 := api.checkTableTags(tableInformation.TableArn, nestorID, t1)
	if err2 != nil {
		t1.Fail(err2)
		return nil, err2
	}
	return tableInformation, nil
}

func (api *DynamoDbAPI) createMonoTable(tableName string, nestorID string, task *reporter.Task) (*TableInformation, error) {
	t0 := task.SubM(reporter.NewMessage("createMonoTable").WithArg("tableName", tableName))

	t1 := t0.Sub("check if table exists")
	tableInformation, err := api.checkTableExistenceAndTags(tableName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}

	if tableInformation != nil {
		t1.Log("table exists")
		t1.Okr(map[string]string{
			"ID":        tableInformation.TableID,
			"arn":       tableInformation.TableArn,
			"tableName": tableInformation.TableName,
		})

		return tableInformation, nil
	}

	t2 := t0.Sub("table does not exist - creating it")
	result, err := api.doCreateMonoTable(tableName, nestorID, t2)
	if err != nil {
		t2.Fail(err)
	}
	t2.Ok()
	t0.Okr(map[string]string{
		"ID":        result.TableID,
		"arn":       result.TableArn,
		"tableName": result.TableName,
	})
	return result, nil
}

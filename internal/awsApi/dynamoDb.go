package awsapi

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nestorhq/nestor/internal/reporter"
)

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

func (api *DynamoDbAPI) createMonoTable(tableName string, task *reporter.Task) {
	t0 := task.SubM(reporter.NewMessage("dynamodb.CreateTableInput").WithArg("tableName", tableName))

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Year"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := api.client.CreateTable(input)
	if err != nil {
		t0.Fail(err)
		os.Exit(1)
	}

	fmt.Println("Created the table", tableName)
}

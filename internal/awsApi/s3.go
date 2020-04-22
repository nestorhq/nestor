package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nestorhq/nestor/internal/reporter"
)

// S3API api
type S3API struct {
	resourceTags *ResourceTags
	client       *s3.S3
}

// S3Information description of a S3
type S3Information struct {
	BucketName string
	BucketArn  string
}

// S3NotificationLambdaDefinition def
type S3NotificationLambdaDefinition struct {
	LambdaArn string
	Prefix    string
	Suffix    string
}

// S3NotificationDefinition def
type S3NotificationDefinition struct {
	Lambdas []S3NotificationLambdaDefinition
}

// MkBucketArn build ARN from bucket name
func MkBucketArn(bucketName string) string {
	return fmt.Sprintf("arn:aws:s3:::%s", bucketName)
}

// NewS3API constructor
func NewS3API(session *session.Session, resourceTags *ResourceTags) (*S3API, error) {
	var api = S3API{resourceTags: resourceTags}
	// Create S3 client
	api.client = s3.New(session)
	return &api, nil
}

func mkS3Information(bucketName string) *S3Information {
	return &S3Information{
		BucketName: bucketName,
		BucketArn:  MkBucketArn(bucketName),
	}
}

func (api *S3API) doTagBucket(bucketName string, nestorID string, t *reporter.Task) error {
	t0 := t.SubM(reporter.NewMessage("api.client.PutBucketTagging").WithArg("bucketName", bucketName))

	tags := api.resourceTags.getTagsAsTagsWithID(nestorID)
	bucketTags := make([]*s3.Tag, 0, 4)
	for _, t := range tags {
		bucketTags = append(bucketTags, &s3.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}

	input := &s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: bucketTags,
		},
	}
	_, err := api.client.PutBucketTagging(input)
	if err != nil {
		t0.Fail(err)
		return err
	}
	return err
}

func (api *S3API) doCreateBucket(bucketName string, nestorID string, t *reporter.Task) (*S3Information, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.CreateBucket").WithArg("bucketName", bucketName))
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err := api.client.CreateBucket(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	return mkS3Information(bucketName), nil
}

func (api *S3API) checkBucketExistenceAndTags(bucketName string, nestorID string, t *reporter.Task) (*S3Information, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.GetBucketTagging").WithArg("bucketName", bucketName))
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}
	result, err := api.client.GetBucketTagging(input)
	if err != nil {
		if getAwsErrorCode(err) == "NoSuchBucket" {
			return nil, nil
		}
		t0.Fail(err)
		return nil, err
	}

	tagsToCheck := map[string]*string{}
	tags := result.TagSet
	for _, tag := range tags {
		tagsToCheck[*tag.Key] = tag.Value
	}
	// check tags
	t1 := t.SubM(reporter.NewMessage("checkTags").WithArgs(tagsToCheck))
	err2 := api.resourceTags.checkTags(tagsToCheck, nestorID)
	if err2 != nil {
		t1.Fail(err2)
		return nil, err2
	}
	t1.Ok()
	return mkS3Information(bucketName), nil
}

func (api *S3API) createBucket(bucketName string, nestorID string, t *reporter.Task) (*S3Information, error) {
	t0 := t.SubM(reporter.NewMessage("checkBucketExistenceAndTags").
		WithArg("bucketName", bucketName).WithArg("id", nestorID))
	result, err := api.checkBucketExistenceAndTags(bucketName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}

	if result != nil {
		t0.Log("s3 bucket exists")
		return result, nil
	}

	t1 := t0.Sub("s3 bucket does not exist - creating it")
	result, err = api.doCreateBucket(bucketName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}
	t1.Ok()

	t2 := t0.Sub("add tags to bucket")
	err = api.doTagBucket(bucketName, nestorID, t2)
	if err != nil {
		t2.Fail(err)
		return nil, err
	}
	t2.Ok()

	t0.Okr(map[string]string{
		"arn":        result.BucketArn,
		"bucketName": result.BucketName,
	})
	return result, nil

}

// GetNotificationConfiguration retrieve bucket notification configuration
func (api *S3API) setNotificationConfiguration(bucketName string, notification *S3NotificationDefinition, t *reporter.Task) error {
	t0 := t.SubM(reporter.NewMessage("api.client.PutBucketNotificationConfiguration").WithArg("bucketName", bucketName))
	input := &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: &s3.NotificationConfiguration{
			LambdaFunctionConfigurations: []*s3.LambdaFunctionConfiguration{},
		},
	}
	var lambdaFunctionConfigurations = input.NotificationConfiguration.LambdaFunctionConfigurations
	for _, lambdaNotif := range notification.Lambdas {
		lambdaFunctionConfigurations = append(lambdaFunctionConfigurations, &s3.LambdaFunctionConfiguration{
			LambdaFunctionArn: &lambdaNotif.LambdaArn,
			Events:            []*string{aws.String("s3:ObjectCreated:*")},
			Filter: &s3.NotificationConfigurationFilter{
				Key: &s3.KeyFilter{
					FilterRules: []*s3.FilterRule{{
						Name:  aws.String("prefix"),
						Value: aws.String(lambdaNotif.Prefix),
					}, {
						Name:  aws.String("suffix"),
						Value: aws.String(lambdaNotif.Suffix),
					}},
				},
			},
		})
	}
	input.NotificationConfiguration.LambdaFunctionConfigurations = lambdaFunctionConfigurations

	_, err := api.client.PutBucketNotificationConfiguration(input)
	if err != nil {
		t0.Fail(err)
		return err
	}
	// t0.LogM(reporter.NewMessage("PutBucketNotificationConfiguration").
	// 	WithArg("input", input.GoString()).
	// 	WithArg("result", result.GoString()))

	return nil

}

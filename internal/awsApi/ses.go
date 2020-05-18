package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/nestorhq/nestor/internal/reporter"
)

// SESAPI api
type SESAPI struct {
	resourceTags *ResourceTags
	region       string
	client       *ses.SES
}

// SesDomainInformation information about a domain
type SesDomainInformation struct {
	SesDomainARN  string
	SesDomainName string
}

// NewSESAPI constructor
func NewSESAPI(session *session.Session, region string, resourceTags *ResourceTags) (*SESAPI, error) {
	var api = SESAPI{
		resourceTags: resourceTags,
		region:       region,
	}
	// Create SES client
	api.client = ses.New(session)
	return &api, nil
}

func (api *SESAPI) createSESDomain(domainName string, task *reporter.Task) (*SesDomainInformation, error) {
	input := ses.ListIdentitiesInput{
		IdentityType: aws.String("Domain"),
		MaxItems:     aws.Int64(64),
	}
	result, err := api.client.ListIdentities(&input)
	if err != nil {
		return nil, nil
	}
	for _, item := range result.Identities {
		if *item == domainName {
			return &SesDomainInformation{
				SesDomainName: domainName,
				SesDomainARN:  fmt.Sprintf("arn:aws:ses:%s:464972470401:identity/%s", api.region, domainName),
			}, nil
		}
	}
	return nil, fmt.Errorf("Feature not implemented: domain %s not created", domainName)
}

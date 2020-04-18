package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/nestorhq/nestor/internal/reporter"
)

// APIGatewayV2API api
type APIGatewayV2API struct {
	resourceTags *ResourceTags
	client       *apigatewayv2.ApiGatewayV2
}

// APIGatewayV2Information description of a ApiGatewayV2
type APIGatewayV2Information struct {
	HTTPApiName     string
	HTTPApiID       string
	HTTPApiEndPoint string
}

func infoAsMap(result *APIGatewayV2Information) map[string]string {
	return map[string]string{
		"apiEndPoint": result.HTTPApiEndPoint,
		"apiID":       result.HTTPApiID,
		"apiName":     result.HTTPApiName,
	}
}

// NewAPIGatewayV2API constructor
func NewAPIGatewayV2API(session *session.Session, resourceTags *ResourceTags) (*APIGatewayV2API, error) {
	var api = APIGatewayV2API{resourceTags: resourceTags}
	// Create ApiGatewayV2 client
	api.client = apigatewayv2.New(session)
	return &api, nil
}
func (api *APIGatewayV2API) getIntegration(apiID string, integrationID string, t *reporter.Task) error {
	t0 := t.SubM(reporter.NewMessage("api.client.GetApi").
		WithArg("apiID", apiID).
		WithArg("integrationID", integrationID))
	input := &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(apiID),
		IntegrationId: aws.String(integrationID),
	}
	result, err := api.client.GetIntegration(input)
	if err != nil {
		t0.Fail(err)
		return err
	}
	fmt.Printf("GetIntegrationInput:%s\n", result.GoString())
	return nil
}

func (api *APIGatewayV2API) getIntegrations(apiID string, t *reporter.Task) error {
	t0 := t.SubM(reporter.NewMessage("api.client.GetIntegrations").WithArg("apiID", apiID))

	var nextToken = ""
	for {
		input := &apigatewayv2.GetIntegrationsInput{
			ApiId:      aws.String(apiID),
			MaxResults: aws.String("32"),
		}
		if len(nextToken) > 0 {
			input.NextToken = &nextToken
		}

		listIntegrations, err := api.client.GetIntegrations(input)
		if err != nil {
			t0.Fail(err)
			return err
		}
		// look for the user pool given by name
		for _, integration := range listIntegrations.Items {
			fmt.Printf("integration: %#v\n", integration)

		}
		// check if we have to paginate
		if listIntegrations.NextToken != nil {
			nextToken = *listIntegrations.NextToken
		} else {
			t0.Ok()
			return nil
		}
	}
}

func (api *APIGatewayV2API) getRoutes(apiID string, t *reporter.Task) error {
	t0 := t.SubM(reporter.NewMessage("api.client.getRoutes").WithArg("apiID", apiID))

	var nextToken = ""
	for {
		input := &apigatewayv2.GetRoutesInput{
			ApiId:      aws.String(apiID),
			MaxResults: aws.String("32"),
		}
		if len(nextToken) > 0 {
			input.NextToken = &nextToken
		}

		listRoutes, err := api.client.GetRoutes(input)
		if err != nil {
			t0.Fail(err)
			return err
		}
		// look for the user pool given by name
		for _, route := range listRoutes.Items {
			fmt.Printf("route: %#v\n", route)

		}
		// check if we have to paginate
		if listRoutes.NextToken != nil {
			nextToken = *listRoutes.NextToken
		} else {
			t0.Ok()
			return nil
		}
	}
}

func (api *APIGatewayV2API) findRestAPIByName(apiName string) (string, error) {
	var nextToken = ""
	for {
		var input = apigatewayv2.GetApisInput{
			MaxResults: aws.String("32"),
		}
		if len(nextToken) > 0 {
			input.NextToken = &nextToken
		}

		listApis, err := api.client.GetApis(&input)
		if err != nil {
			return "", err
		}
		// look for the user pool given by name
		for _, api := range listApis.Items {
			if *api.Name == apiName {
				return *api.ApiId, nil
			}
		}
		// check if we have to paginate
		if listApis.NextToken != nil {
			nextToken = *listApis.NextToken
		} else {
			// the pool was not found
			return "", nil
		}
	}
}

func (api *APIGatewayV2API) getAPIByID(apiID string, nestorID string, t *reporter.Task) (*APIGatewayV2Information, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.GetApi").WithArg("apiID", apiID))
	input := &apigatewayv2.GetApiInput{
		ApiId: aws.String(apiID),
	}
	result, err := api.client.GetApi(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	// check tags
	t1 := t.SubM(reporter.NewMessage("checkTags").WithArgs(result.Tags))
	err = api.resourceTags.checkTags(result.Tags, nestorID)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}
	t1.Ok()

	return &APIGatewayV2Information{
		HTTPApiID:       *result.ApiId,
		HTTPApiName:     *result.Name,
		HTTPApiEndPoint: *result.ApiEndpoint,
	}, nil
}

func (api *APIGatewayV2API) doCreateRestAPI(apiName string, lambdaTargetArn string, nestorID string, t *reporter.Task) (*APIGatewayV2Information, error) {
	t0 := t.SubM(reporter.NewMessage("api.client.CreateApi").WithArg("apiName", apiName))
	input := &apigatewayv2.CreateApiInput{
		Name:         aws.String(apiName),
		Tags:         aws.StringMap(api.resourceTags.getTagsAsMapWithID(nestorID)),
		ProtocolType: aws.String("HTTP"),
		Target:       aws.String(lambdaTargetArn),
	}
	result, err := api.client.CreateApi(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	return &APIGatewayV2Information{
		HTTPApiID:       *result.ApiId,
		HTTPApiName:     *result.Name,
		HTTPApiEndPoint: *result.ApiEndpoint,
	}, nil

}

func (api *APIGatewayV2API) checkRestAPIExistenceAndTags(apiName string, nestorID string, t *reporter.Task) (*APIGatewayV2Information, error) {
	t0 := t.SubM(reporter.NewMessage("findRestAPIByName").WithArg("apiName", apiName))
	apiID, err := api.findRestAPIByName(apiName)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if apiID == "" {
		t0.Log("REST API not found")
		return nil, nil
	}

	t1 := t.SubM(reporter.NewMessage("getAPIByID").WithArg("apiID", apiID))
	result, err := api.getAPIByID(apiID, nestorID, t1)

	return result, nil
}

func (api *APIGatewayV2API) createRestAPI(apiName string, lambdaTargetArn string, nestorID string, t *reporter.Task) (*APIGatewayV2Information, error) {
	t0 := t.SubM(reporter.NewMessage("checkRestAPIExistenceAndTags").
		WithArg("apiName", apiName).WithArg("nestorID", nestorID))
	result, err := api.checkRestAPIExistenceAndTags(apiName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if result != nil {
		t0.Log("rest API exists")
		t0.Okr(infoAsMap(result))
		return result, nil
	}

	t1 := t0.Sub("rest API does not exist - creating it")
	result, err = api.doCreateRestAPI(apiName, lambdaTargetArn, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}
	t1.Ok()
	t0.Okr(infoAsMap(result))
	return result, nil
}

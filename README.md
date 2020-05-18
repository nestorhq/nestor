# nestor
Nestor CLI

# Templates
https://github.com/rakyll/statik

command (to be verified)

```
$ statik -src templates -dest internal/templates -f
```

# Interesting
https://github.com/retgits/lambda-util/blob/master/s3.go

# TODO
- add command to display list of allocated resources with arn, names, etc...
- allow variables in environment variables
- add duration for lambda
- fix in 	`internal/awsapi/lambda.go` hard coded region:
```
sourceArn := fmt.Sprintf("arn:aws:execute-api:us-west-2:%s:%s/*/$default", api.account, apiID)
```
- allow change of throughput for dynamodb tables
- add creation of SES domains


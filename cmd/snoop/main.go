package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(context context.Context, e interface{}) (string, error) {
	out, err := json.Marshal(e)
	if err != nil {
		return "error decoding", err
	}

	fmt.Println(string(out))
	return fmt.Sprintf(">lambda>go>snoop> %s", string(out)), nil
}
func main() {
	lambda.Start(handleRequest)
}

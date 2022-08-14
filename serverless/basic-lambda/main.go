package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LambdaEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type LambdaResponse struct {
	Message string `json:"msg"`
}

func handleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	fmt.Printf("Running %s", lc.ClientContext.Client.AppPackageName)
	return LambdaResponse{Message: fmt.Sprintf("hello %s , Your age is %d Goodby !!!", event.Name, event.Age)}, nil
}

func main() {
	lambda.Start(handleRequest)
}

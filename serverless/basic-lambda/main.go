package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v4"
)

// Commenting this to model request and response based on API Gateway Proxy
//type LambdaEvent struct {
//	Name string `json:"name"`
//	Age  int    `json:"age"`
//}

// type LambdaResponse struct {

// 	Message string `json:"msg"`
// }

// func handleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
// 	lc, _ := lambdacontext.FromContext(ctx)
// 	fmt.Printf("Running %s", lc.ClientContext.Client.AppPackageName)
// 	return LambdaResponse{Message: fmt.Sprintf("hello %s , Your age is %d Goodby !!!", event.Name, event.Age)}, nil
// }

type AWSCognitoClaims struct {
	Client_ID string `json:client_id`
	Username  string `json:username`
	jwt.StandardClaims
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//event := LambdaEvent{}
	//error := json.Unmarshal([]byte(req.Body), &event)
	tokenString := req.Headers["Authorization"]
	if len(tokenString) <= 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unable to find Auth Token",
		}, nil

	}
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unable to extract claim from Token",
		}, nil
	}
	name := claims["cognito:username"].(string)

	// Getting name from Cognito Token now
	//name := req.QueryStringParameters["name"]
	if len(name) <= 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Unable to find your name",
		}, nil
	}
	age, err := strconv.Atoi(req.QueryStringParameters["age"])

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Unable to parse your age.",
		}, nil
	}

	if age <= 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Your Age cant be negative.",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("hello %s , Your age is %d", name, age),
	}, nil

}

func main() {
	lambda.Start(handleRequest)
}

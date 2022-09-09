package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Wahts should be your Userame ?")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")
	fmt.Println("Wahts should be your Password ?")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSuffix(password, "\n")
	fmt.Println("Wahts your Cognito Client ID ?")
	clientid, _ := reader.ReadString('\n')
	clientid = strings.TrimSuffix(clientid, "\n")
	fmt.Println("Wahts your API Gateway URL ?")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSuffix(url, "\n")
	fmt.Println("Starting AWS Connection for SignUp....")
	conf := &aws.Config{Region: aws.String("us-east-1")}
	session, _ := session.NewSession(conf)
	cognitoSession := cognito.New(session)
	authInput := &cognito.InitiateAuthInput{
		ClientId: aws.String(clientid),
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	auth, err := cognitoSession.InitiateAuth(authInput)

	if err != nil {
		fmt.Println("Authentication Initiation failed-")
		fmt.Println(err)
		return
	}

	token := auth.AuthenticationResult.IdToken
	fmt.Println("Calling with Age query string. Expect API gateway to execute Lambda and provide response.")
	// Success Scenario
	requestURL := fmt.Sprintf("%s?age=32", url)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return
	}
	req.Header.Set("Authorization", *token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	//Negative Scenario
	fmt.Println("Calling without Age query string. Expect API gateway validation error.")
	requestURL = url
	req, err = http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return
	}
	req.Header.Set("Authorization", *token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

}

.DEFAULT_GOAL := basiclambdadeploy
basiclambdabuild:
	@echo "Building Basic Lambda function"
	cd basic-lambda && go mod tidy
	cd basic-lambda && GOOS=linux GOARCH=amd64 go build -o main main.go

basiclambdadeploy:
	@echo "Deploying to AWS using Terraform Cloud"
	cdktf deploy
package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

const projectname = "goserverlessaws"
const version = "1.0"

// PolicyDocument is our definition of our policies to be uploaded to IAM.
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type Principal struct {
	Service string
}

// StatementEntry will dictate what this policy will allow or not allow.
type StatementEntry struct {
	Effect    string
	Action    []string
	Principal Principal
}

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	// The code that defines your stack goes here
	aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-east-1"),
		DefaultTags: &aws.AwsProviderDefaultTags{
			Tags: &map[string]*string{
				"CreatedBy":   jsii.String("cdktf"),
				"Environment": jsii.String("dev"),
			},
		},
	})

	asset := cdktf.NewTerraformAsset(stack, jsii.String("lambda-asset"), &cdktf.TerraformAssetConfig{
		Path: jsii.String("basic-lambda"),
		Type: cdktf.AssetType_ARCHIVE,
	})

	// Creating S3 bucket for keeping Lambda Artifacts
	bucket := s3.NewS3Bucket(stack, jsii.String("s3lambdabucket"), &s3.S3BucketConfig{
		BucketPrefix: jsii.String(projectname),
	})

	lambdaArchivesObj := s3.NewS3BucketObject(stack, jsii.String(fmt.Sprintf("%s-basiclambda", projectname)), &s3.S3BucketObjectConfig{
		Bucket: bucket.Bucket(),
		Key:    jsii.String(fmt.Sprintf("%s-%s", version, *asset.FileName())),
		Source: asset.Path(),
	})
	assumePolicy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"sts:AssumeRole",
				},
				Principal: Principal{
					Service: "lambda.amazonaws.com",
				},
			},
		},
	}
	assumeRolePolicy, err := json.Marshal(&assumePolicy)
	if err != nil {
		fmt.Println("Error in marshaling the policy", err)
		return nil
	}

	lambdarole := iam.NewIamRole(stack, jsii.String("lambdaiam"), &iam.IamRoleConfig{
		NamePrefix:       jsii.String(projectname),
		AssumeRolePolicy: jsii.String(string(assumeRolePolicy)),
	})

	iam.NewIamRolePolicyAttachment(stack, jsii.String("lambda-managed-policy"), &iam.IamRolePolicyAttachmentConfig{
		Role:      lambdarole.Name(),
		PolicyArn: jsii.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})
	basiclambda := lambdafunction.NewLambdaFunction(stack, jsii.String("basic-lambda"), &lambdafunction.LambdaFunctionConfig{
		FunctionName:   jsii.String(fmt.Sprintf("%s-%s", projectname, "basiclambda")),
		S3Bucket:       bucket.Bucket(),
		S3Key:          lambdaArchivesObj.Key(),
		Handler:        jsii.String("main"),
		Runtime:        jsii.String("go1.x"),
		Role:           lambdarole.Arn(),
		SourceCodeHash: cdktf.Fn_Filebase64sha256(lambdaArchivesObj.Source()),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("basic-lamga-arn"), &cdktf.TerraformOutputConfig{
		Value: basiclambda.Arn(),
	})
	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "serverless")
	cdktf.NewRemoteBackend(stack, &cdktf.RemoteBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("erankitcs"),
		Workspaces:   cdktf.NewNamedRemoteWorkspace(jsii.String("goserverlesaws")),
	})

	app.Synth()
}

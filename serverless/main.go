package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

const projectname = "goserverlessaws"
const version = "1.0"

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

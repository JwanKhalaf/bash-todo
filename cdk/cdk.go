package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/constructs-go/constructs/v10"
)

type GoTodoAppStackProps struct {
	awscdk.StackProps
}

func NewGoTodoAppStack(scope constructs.Construct, id string, props *GoTodoAppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("Todo"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("task"),
			Type: awsdynamodb.AttributeType_STRING},
	})

	bundlingOptions := &awscdklambdagoalpha.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w" -tags lambda.norpc`)},
	}

	handler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Function"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/getitems/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	table.GrantReadWriteData(handler)

	integration := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("mylambdaintegration"), handler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
		PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
	})

	api := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("Api"), &awscdkapigatewayv2alpha.HttpApiProps{
		DefaultIntegration: integration,
	})

	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{Value: api.Url()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewGoTodoAppStack(app, "GoTodoAppStack", &GoTodoAppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

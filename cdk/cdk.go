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

	// create a new stack
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create a dynamodb table to store the tasks
	table := awsdynamodb.NewTable(stack, jsii.String("tasks"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("task_id"),
			Type: awsdynamodb.AttributeType_STRING},
		BillingMode:         awsdynamodb.BillingMode_PAY_PER_REQUEST,
		TimeToLiveAttribute: jsii.String("time_to_live"),
	})

	// add a global secondary index based on user_id
	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName:    jsii.String("user-index"),
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("user_id"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:      &awsdynamodb.Attribute{Name: jsii.String("created_at"), Type: awsdynamodb.AttributeType_STRING},
	})

	// bundling options to make go fast
	bundlingOptions := &awscdklambdagoalpha.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w" -tags lambda.norpc`)},
	}

	// creating the aws lambda for listing items
	listItemsHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("ListItemsFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/items/list/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	// grant dynamodb read write permissions to the list items lambda
	table.GrantReadWriteData(listItemsHandler)

	// creating the aws lambda for creating a new item
	createItemHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("CreateItemFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../api/items/create/lambda"),
		Environment:  &map[string]*string{"DYNAMODB_TABLENAME": table.TableName()},
		Bundling:     bundlingOptions,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Millis(jsii.Number(15000)),
	})

	// grant dynamodb read write permissions to the create item lambda
	table.GrantReadWriteData(createItemHandler)

	// create a new http itemsApi gateway
	itemsApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("ItemsApi"), &awscdkapigatewayv2alpha.HttpApiProps{})

	// add route for listing the tasks
	itemsApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/tasks"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("listItemsLambdaIntegration"), listItemsHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// add route for creating a new task
	itemsApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/tasks"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("createItemLambdaIntegration"), createItemHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// output the lambda url to the console
	awscdk.NewCfnOutput(stack, jsii.String("ItemsApiUrl"), &awscdk.CfnOutputProps{Value: itemsApi.Url()})

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

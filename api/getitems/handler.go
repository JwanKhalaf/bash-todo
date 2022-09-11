package getitems

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"net/http"
	"os"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	const dynamodbTableNameEnvKey = "DYNAMODB_TABLENAME"

	dynamodbTableName, ok := os.LookupEnv(dynamodbTableNameEnvKey)

	if !ok {
		log.Fatalf("the %v variable was not set!", dynamodbTableNameEnvKey)
	} else {
		log.Printf("The %v variable is set to: %v", dynamodbTableNameEnvKey, dynamodbTableName)
	}

	log.Println("Running the GetItemsHandler!")

	// using the sdk's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))

	if err != nil {
		log.Fatalf("unable to load sdk config, %v", err)
	}

	// using the config value, create the dynamodb client
	dynamodbService := dynamodb.NewFromConfig(cfg)

	var tasks []Task

	// build the required scan params
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(dynamodbTableName),
	}

	response, err := dynamodbService.Scan(context.TODO(), scanInput)

	if err != nil {
		log.Printf("could not scan the dyanmodb table! error: %v", err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &tasks)

		if err != nil {
			log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(tasks)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResponse)
}

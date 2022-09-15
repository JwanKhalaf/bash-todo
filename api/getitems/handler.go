package getitems

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"net/http"
	"os"
)

var dynamodbTableName string
var dynamodbService *dynamodb.Client

func init() {
	dynamodbTableName, ok := os.LookupEnv("DYNAMODB_TABLENAME")
	if !ok {
		log.Fatal("the DYNAMODB_TABLENAME variable was not set!")
	}

	log.Printf("The DYNAMODB_TABLENAME variable is set to: %v", dynamodbTableName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load sdk config: %v", err)
	}

	dynamodbService = dynamodb.NewFromConfig(cfg)
}

func listTasks(ctx context.Context) (tasks []Task, err error) {
	result, err := dynamodbService.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(dynamodbTableName),
	})
	if err != nil {
		err = fmt.Errorf("could not scan the dyanmodb table: %w", err)
		return
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &tasks)
	return
}

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GetItemsHandler!")
	w.Header().Set("Content-Type", "application/json")

	tasks, err := listTasks(r.Context())
	if err != nil {
		log.Printf("GetItemsHandler: failed to list tasks: %v", err)
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Printf("GetItemsHandler: error in JSON marshal: %v", err)
	}
}

package getitems

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type taskStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewTaskStore() *taskStore {
	dynamodbTableName, ok := os.LookupEnv("DYNAMODB_TABLENAME")
	if !ok {
		log.Fatal("the DYNAMODB_TABLENAME variable was not set!")
	}

	log.Printf("The DYNAMODB_TABLENAME variable is set to: %v", dynamodbTableName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load sdk config: %v", err)
	}

	return &taskStore{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: dynamodbTableName,
	}
}

func (d *taskStore) ListTasks(ctx context.Context) ([]Task, error) {
	response, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("could not scan the dyanmodb table: %w", err)
	}

	var tasks []Task
	err = attributevalue.UnmarshalListOfMaps(response.Items, &tasks)
	return tasks, err
}

type TaskLister interface {
	ListTasks(context.Context) ([]Task, error)
}

func GetItemsHandler(lister TaskLister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running the GetItemsHandler!")

		w.Header().Set("Content-Type", "application/json")

		tasks, err := lister.ListTasks(r.Context())
		if err != nil {
			log.Printf("GetItemsHandler: failed to list tasks: %v", err)
			http.Error(w, "failed to list tasks", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			log.Printf("GetItemsHandler: error in JSON marshal: %v", err)
		}
	})
}

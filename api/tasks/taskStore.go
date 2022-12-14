package tasks

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type TaskStore struct {
	client    *dynamodb.Client
	tableName string
}

type TasksRepository interface {
	GetTask(ctx context.Context, taskID string) (Task, error)
	ListTasks(context.Context) ([]Task, error)
	CreateTask(ctx context.Context, task string) (string, error)
	UpdateTask(ctx context.Context, task Task) error
}

func NewTaskStore() *TaskStore {
	dynamodbTableName, ok := os.LookupEnv("DYNAMODB_TABLENAME")
	if !ok {
		log.Fatal("the DYNAMODB_TABLENAME variable was not set!")
	}

	log.Printf("The DYNAMODB_TABLENAME variable is set to: %v", dynamodbTableName)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load sdk config: %v", err)
	}

	return &TaskStore{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: dynamodbTableName,
	}
}

func (d *TaskStore) GetTask(ctx context.Context, taskID string) (Task, error) {
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"task_id": &types.AttributeValueMemberS{Value: taskID},
		},
	})
	if err != nil {
		return Task{}, fmt.Errorf("could not get the item from the dyanmodb table: %w", err)
	}

	var task Task
	err = attributevalue.UnmarshalMap(response.Item, &task)
	return task, err
}

func (d *TaskStore) ListTasks(ctx context.Context) ([]Task, error) {
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

func (d *TaskStore) CreateTask(ctx context.Context, task string) (string, error) {
	staticUserId := "8600aab6-d540-4228-8af5-35218bd564a6"

	item := Task{
		TaskID:    uuid.New().String(),
		UserID:    staticUserId,
		Content:   task,
		CreatedAt: time.Now().Format(time.RFC3339),
		IsDone:    false,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return "", fmt.Errorf("could not marshal map task: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	})
	if err != nil {
		return "", fmt.Errorf("could not put item into dynamodb table: %w", err)
	}

	return item.TaskID, err
}

func (d *TaskStore) UpdateTask(ctx context.Context, task Task) error {
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}

	update := expression.Set(expression.Name("content"), expression.Value(task.Content))
	update.Set(expression.Name("is_done"), expression.Value(task.IsDone))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("could not build expression for update, here's why: %v\n", err)
	} else {
		response, err = d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(d.tableName),
			Key: map[string]types.AttributeValue{
				"task_id": &types.AttributeValueMemberS{Value: task.TaskID},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		})
		if err != nil {
			log.Printf("coult not update task %v, here's why: %v\n", task.TaskID, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("could not unmarshall update response, here's why: %v\n", err)
			}
		}
	}

	return err
}

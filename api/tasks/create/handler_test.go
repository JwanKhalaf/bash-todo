package create

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jwankhalaf/bash-todo/api/tasks"
)

type StubTaskStore struct {
	getTask    func(ctx context.Context, taskID string) (tasks.Task, error)
	listTasks  func(ctx context.Context) ([]tasks.Task, error)
	createTask func(ctx context.Context, taskContent string) (string, error)
	updateTask func(ctx context.Context, task tasks.Task) error
}

func (s *StubTaskStore) GetTask(ctx context.Context, taskID string) (tasks.Task, error) {
	return s.getTask(ctx, taskID)
}

func (s *StubTaskStore) ListTasks(ctx context.Context) ([]tasks.Task, error) {
	return s.listTasks(ctx)
}

func (s *StubTaskStore) CreateTask(ctx context.Context, taskContent string) (string, error) {
	return s.createTask(ctx, taskContent)
}

func (s *StubTaskStore) UpdateTask(ctx context.Context, task tasks.Task) error {
	return s.updateTask(ctx, task)
}

func TestCreateTaskHandler(t *testing.T) {
	// create the stub task store
	expectedTaskID := CreateTaskResponse{TaskID: "test_task_id"}
	expectedTaskContent := "expected_task_content"

	taskStore := StubTaskStore{
		createTask: func(ctx context.Context, taskContent string) (string, error) {
			if taskContent != expectedTaskContent {
				t.Errorf("got: CreateTask(%s) expected CreateTask(%s)", taskContent, expectedTaskContent)
			}
			return expectedTaskID.TaskID, nil
		},
	}

	jsonValue, _ := json.Marshal(CreateTaskRequest{Content: expectedTaskContent})

	// create a request to pass to the handler
	request, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))

	// set the content type
	request.Header.Set("content-type", "application/json")

	// create a response recorder
	response := httptest.NewRecorder()

	handler := GetCreateTaskHandler(&taskStore)

	handler.ServeHTTP(response, request)

	// decode the json response into TaskID
	got := getTaskIDFromResponse(t, response.Body)

	// check the status code is what we expect
	assertStatusCode(t, response.Code, http.StatusCreated)

	// check the response body is what we expect
	assertTaskID(t, got, expectedTaskID)
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func assertTaskID(t testing.TB, got, want CreateTaskResponse) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}

func getTaskIDFromResponse(t testing.TB, body io.Reader) (taskID CreateTaskResponse) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&taskID)

	if err != nil {
		t.Fatalf("unable to process response from server %q into CreateTaskResponse, '%v'", body, err)
	}

	return
}

package list

import (
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
	tasks []tasks.Task
}

func (s *StubTaskStore) ListTasks(ctx context.Context) ([]tasks.Task, error) {
	return s.tasks, nil
}

func (s *StubTaskStore) CreateTask(ctx context.Context, taskContent string) (string, error) {
	task := tasks.Task{TaskID: "test_task_id", UserID: "test_user_id", Content: taskContent, CreatedAt: "test_created_at", IsDone: false}
	s.tasks = append(s.tasks, task)
	return task.TaskID, nil
}

func TestListTasksHandler(t *testing.T) {
	// create the stub task store
	taskStore := StubTaskStore{
		tasks: []tasks.Task{
			{TaskID: "test_task_id", UserID: "test_user_id", Content: "test_content", CreatedAt: "test_created_at", IsDone: false},
		},
	}

	// create a request to pass to our handler
	request, _ := http.NewRequest("GET", "/", nil)

	// we create a response recorder which satisfies http.responsewriter to
	// to record the response
	response := httptest.NewRecorder()

	handler := GetListTasksHandler(&taskStore)

	// our handler satisfies http.handler, so we can call its serve http method
	// directly and pass in our request and response recorder
	handler.ServeHTTP(response, request)

	// decode the json response into []tasks.Task
	got := getTasksFromResponse(t, response.Body)

	// check the status code is what we expect
	assertStatusCode(t, response.Code, http.StatusOK)

	// check the response body is what we expect
	assertTasks(t, got, taskStore.tasks)
}

func getTasksFromResponse(t testing.TB, body io.Reader) (tasks []tasks.Task) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		t.Fatalf("unable to process response from server %q into slice of Task, '%v'", body, err)
	}

	return
}

func assertTasks(t testing.TB, got, want []tasks.Task) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

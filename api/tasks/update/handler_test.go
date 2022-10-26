package update

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestUpdateTask(t *testing.T) {
	t.Run("make sure the updated task is what gets processed and that a status code of no content is returned", func(t *testing.T) {
		// expected updated task
		expectedUpdatedTask := tasks.Task{TaskID: "test_task_id", UserID: "test_user_id", Content: "test_task_content", CreatedAt: "test_created_at", IsDone: true}

		// create the stub task store
		taskStore := StubTaskStore{
			updateTask: func(ctx context.Context, task tasks.Task) error {
				if diff := cmp.Diff(task, expectedUpdatedTask); diff != "" {
					t.Errorf("got: UpdateTask(%v) expected UpdateTask(%v)", task, expectedUpdatedTask)
				}
				return nil
			},
		}

		jsonValue, _ := json.Marshal(expectedUpdatedTask)

		// create a request to pass to the handler
		request, _ := http.NewRequest("PUT", "/tasks/test_task_id", bytes.NewBuffer(jsonValue))

		// set the content type
		request.Header.Set("content-type", "application/json")

		// create a response recorder
		response := httptest.NewRecorder()

		handler := GetUpdateTaskHandler(&taskStore)

		handler.ServeHTTP(response, request)

		// check the status code is what we expect
		assertStatusCode(t, response.Code, http.StatusNoContent)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

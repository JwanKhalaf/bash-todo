package get

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func TestGetTaskHandler(t *testing.T) {
	t.Run("get 404 when task does not exist", func(t *testing.T) {
		expectedTaskID := "test_task_id"
		taskStore := StubTaskStore{
			getTask: func(ctx context.Context, taskID string) (tasks.Task, error) {
				if taskID != expectedTaskID {
					t.Errorf("got: GetTask(ctx, %s) expected GetTask(ctx, %s)", taskID, expectedTaskID)
				}
				return tasks.Task{}, errors.New("no task found")
			},
		}

		request, _ := http.NewRequest("GET", fmt.Sprintf("/tasks/%v", expectedTaskID), nil)

		response := httptest.NewRecorder()

		handler := GetGetTaskHandler(&taskStore)

		handler.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("get 202 and the task when task does exist", func(t *testing.T) {
		expectedTaskID := "test_task_id"
		expectedTask := tasks.Task{TaskID: expectedTaskID, UserID: "test_user_id", Content: "test_content", CreatedAt: "test_created_at", IsDone: false}

		taskStore := StubTaskStore{
			getTask: func(ctx context.Context, taskID string) (tasks.Task, error) {
				if taskID != expectedTaskID {
					t.Errorf("got: GetTask(ctx, %s) expected GetTask(ctx, %s)", taskID, expectedTaskID)
				}
				return expectedTask, nil
			},
		}

		request, _ := http.NewRequest("GET", fmt.Sprintf("/tasks/%v", expectedTaskID), nil)

		response := httptest.NewRecorder()

		handler := GetGetTaskHandler(&taskStore)

		handler.ServeHTTP(response, request)

		got := getTaskFromResponse(t, response.Body)

		assertStatusCode(t, response.Code, http.StatusOK)

		assertTask(t, got, expectedTask)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func assertTask(t testing.TB, got, want tasks.Task) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}

func getTaskFromResponse(t testing.TB, body io.Reader) (task tasks.Task) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&task)

	if err != nil {
		t.Fatalf("unable to process response from server %q into Task, '%v'", body, err)
	}

	return
}

package list

import (
	"net/http"
	"net/http/httptest"
	"testing"

    "github.com/jwankhalaf/bash-todo/api/tasks"
)

func TestListTasksHandler(t *testing.T) {
  // create the task store
  taskStore := tasks.NewTaskStore()

  // create a request to pass to our handler
  request, err := http.NewRequest("POST", "/", nil)
  if err != nil {
    t.Fatal(err)
  }

  // we create a response recorder which satisfies http.responsewriter to
  // to record the response
  response := httptest.NewRecorder()

  handler := GetListTasksHandler(taskStore)

  // our handler satisfies http.handler, so we can call its serve http method
  // directly and pass in our request and response recorder
  handler.ServeHTTP(response, request)

  // check the status code is what we expect
  if status := response.Code; status != http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
  }

  // check the response body is what we expect
  expected := `[]`
  if response.Body.String() != expected {
    t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
  }
}

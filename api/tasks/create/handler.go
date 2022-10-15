package create

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/jwankhalaf/bash-todo/api/tasks"
)

func GetCreateTaskHandler(repository tasks.TasksRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running the CreateTaskHandler!")

		type CreateTaskRequest struct {
			Content string `json:"content"`
		}

		type CreateTaskResponse struct {
			TaskID string `json:"task_id"`
		}

		// enforce a json content-type
		contentType := r.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if mediatype != "application/json" {
			http.Error(w, "api expects application/json content-type", http.StatusUnsupportedMediaType)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var createRequestTask CreateTaskRequest
		if err := dec.Decode(&createRequestTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		taskID, err := repository.CreateTask(r.Context(), createRequestTask.Content)
		if err != nil {
			log.Printf("CreateTaskHandler: failed to create task: %v", err)
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(CreateTaskResponse{TaskID: taskID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

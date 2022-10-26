package update

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/jwankhalaf/bash-todo/api/tasks"
)

func GetUpdateTaskHandler(repository tasks.TasksRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("running the update task handler...")

		// enforce a json content-type
		contentType := r.Header.Get("content-type")
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
		var task tasks.Task
		if err := dec.Decode(&task); err != nil {
			http.Error(w, "request body is invalid", http.StatusBadRequest)
			return
		}

		err = repository.UpdateTask(r.Context(), task)
		if err != nil {
			log.Printf("UpdateTaskHandler: failed to update task: %v", err)
			http.Error(w, "failed to update task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

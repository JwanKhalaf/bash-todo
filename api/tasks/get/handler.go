package get

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jwankhalaf/bash-todo/api/tasks"
)

func GetGetTaskHandler(repository tasks.TasksRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("running the get task handler...")
		log.Printf("the request url is: %q", r.URL)

		// get the url param `:task-id` in `/tasks/:task-id`
		taskID := strings.TrimPrefix(r.URL.Path, "/tasks/")

		w.Header().Set("content-type", "application/json")

		task, err := repository.GetTask(r.Context(), taskID)
		if err != nil {
			log.Printf("GetTaskHandler: failed to get task: %v", err)
			http.Error(w, "failed to get task", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			log.Printf("GetTaskHandler: error in JSON marshal: %v", err)
		}
	})
}

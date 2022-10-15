package create

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jwankhalaf/bash-todo/api/tasks"
)

func GetCreateItemHandler(repository tasks.TasksRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running the CreateItemHandler!")

		w.Header().Set("Content-Type", "application/json")

		taskDescription := "Deploy the application"

		// err := json.NewDecoder(r.Body).Decode(taskDescription)
		// if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		// return
		// }

		taskID, err := repository.CreateTask(r.Context(), taskDescription)
		if err != nil {
			log.Printf("CreateItemHandler: failed to create task: %v", err)
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(taskID)
		if err != nil {
			log.Printf("GetItemsHandler: error in JSON marshal: %v", err)
		}
	})
}

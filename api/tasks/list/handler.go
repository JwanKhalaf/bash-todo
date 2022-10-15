package list

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jwankhalaf/bash-todo/api/items"
)

func GetListItemsHandler(repository items.TasksRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running the ListItemsHandler!")

		w.Header().Set("Content-Type", "application/json")

		tasks, err := repository.ListTasks(r.Context())
		if err != nil {
			log.Printf("GetItemsHandler: failed to list tasks: %v", err)
			http.Error(w, "failed to list tasks", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			log.Printf("GetItemsHandler: error in JSON marshal: %v", err)
		}
	})
}

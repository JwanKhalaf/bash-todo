package getitems

import (
	"log"
	"net/http"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GetItemsHandler!")

	w.Write([]byte("Hello from the getItems handler!"))
}

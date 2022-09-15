package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/getitems"
)

func main() {
	log.Println("Running the main function!")

	mux := http.NewServeMux()

	mux.Handle("/", getitems.GetItemsHandler(getitems.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

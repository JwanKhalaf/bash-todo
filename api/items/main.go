package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/items/list"
	"github.com/jwankhalaf/bash-todo/api/items/create"
)

func main() {
	log.Println("Running the main function!")

	mux := http.NewServeMux()

	mux.Handle("/", list.GetListItemsHandler(NewTaskStore()))
	mux.Handle("/", create.GetCreateItemHandler(NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

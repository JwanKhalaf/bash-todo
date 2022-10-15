package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/tasks"
	"github.com/jwankhalaf/bash-todo/api/tasks/list"
)

func main() {
	log.Println("running the list tasks lambda!")

	mux := http.NewServeMux()

	mux.Handle("/", list.GetListItemsHandler(tasks.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

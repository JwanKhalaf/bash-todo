package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/tasks"
	"github.com/jwankhalaf/bash-todo/api/tasks/create"
)

func main() {
	log.Println("running the create task lambda!")

	mux := http.NewServeMux()

	mux.Handle("/", create.GetCreateTaskHandler(tasks.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

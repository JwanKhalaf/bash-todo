package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/tasks"
	"github.com/jwankhalaf/bash-todo/api/tasks/update"
)

func main() {
	log.Println("running the update task lambda...")

	mux := http.NewServeMux()

	mux.Handle("/", update.GetUpdateTaskHandler(tasks.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

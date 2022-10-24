package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/tasks"
	"github.com/jwankhalaf/bash-todo/api/tasks/get"
)

func main() {
	log.Println("running the get task lambda...")

	mux := http.NewServeMux()

	mux.Handle("/", get.GetGetTaskHandler(tasks.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

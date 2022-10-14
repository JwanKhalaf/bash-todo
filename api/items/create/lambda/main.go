package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/items"
	"github.com/jwankhalaf/bash-todo/api/items/create"
)

func main() {
	log.Println("running the create task lambda!")

	mux := http.NewServeMux()

	mux.Handle("/", create.GetCreateItemHandler(items.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

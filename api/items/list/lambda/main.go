package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/bash-todo/api/items"
	"github.com/jwankhalaf/bash-todo/api/items/list"
)

func main() {
	log.Println("running the list tasks lambda!")

	mux := http.NewServeMux()

	mux.Handle("/", list.GetListItemsHandler(items.NewTaskStore()))
	algnhsa.ListenAndServe(mux, nil)
}

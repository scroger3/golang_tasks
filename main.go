package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Product struct
type Product struct {
	ID     int    `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}

// Error struct
type Error struct {
	Error string `json:"Error"`
}

//Items - slice of products
var Items []Product

func notFound(w http.ResponseWriter, r *http.Request, e Error) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(e)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	if len(Items) < 1 {
		notFound(w, r, Error{Error: "No one item exists"})
		return
	}

	json.NewEncoder(w).Encode(Items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	found := false
	vars := mux.Vars(r)
	e := Error{Error: "This item doesn't exist"}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		notFound(w, r, e)

		return
	}

	for _, item := range Items {
		if item.ID == id {
			found = true

			json.NewEncoder(w).Encode(item)

			break
		}
	}

	if !found {
		notFound(w, r, e)
	}
}

func addItem(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var item Product

	json.Unmarshal(reqBody, &item)
	Items = append(Items, item)

	w.WriteHeader(http.StatusCreated)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	found := false
	vars := mux.Vars(r)
	e := Error{Error: "Item with this id doesn't exist"}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		notFound(w, r, e)

		return
	}

	for index, item := range Items {
		if item.ID == id {
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &Items[index])
			found = true

			w.WriteHeader(http.StatusAccepted)
			break
		}
	}

	if !found {
		notFound(w, r, e)
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	found := false
	vars := mux.Vars(r)
	e := Error{Error: "Item with this id doen't exist"}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		notFound(w, r, e)

		return
	}

	for index, item := range Items {
		if item.ID == id {
			Items = append(Items[:index], Items[index+1:]...)
			found = true

			w.WriteHeader(http.StatusAccepted)
			break
		}
	}

	if !found {
		notFound(w, r, e)
	}
}

func main() {
	Items = []Product{}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/item/{id:[0-9]+}", getItem).Methods("GET")
	r.HandleFunc("/item", addItem).Methods("POST")
	r.HandleFunc("/item/{id:[0-9]+}", updateItem).Methods("PUT")
	r.HandleFunc("/item/{id:[0-9]+}", deleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

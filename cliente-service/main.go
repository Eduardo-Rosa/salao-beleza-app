package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Cliente struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

var clientes []Cliente
var currentID int = 1

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/clientes", getClientes).Methods("GET")
	r.HandleFunc("/clientes", createCliente).Methods("POST")
	r.HandleFunc("/clientes/{id}", getCliente).Methods("GET")

	log.Println("Servidor rodando na porta 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func getClientes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func createCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	_ = json.NewDecoder(r.Body).Decode(&cliente)
	cliente.ID = currentID
	currentID++
	clientes = append(clientes, cliente)
	w.Header().Set("Content.Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

func getCliente(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, cliente := range clientes {
		if cliente.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cliente)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(nil)
}

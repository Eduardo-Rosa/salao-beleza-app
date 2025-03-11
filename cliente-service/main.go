package main

import (
	"database/sql"
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

var db *sql.DB

func main() {

	var err error
	connStr := "user=root dbname=clientes_db sslmode=disable password=pass"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/clientes", getClientes).Methods("GET")
	r.HandleFunc("/clientes", createCliente).Methods("POST")
	r.HandleFunc("/clientes/{id}", getCliente).Methods("GET")

	log.Println("Servidor rodando na porta 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func getClientes(w http.ResponseWriter, r *http.Request) {
	var clientes []Cliente

	rows, err := db.Query("SELECT id, nome, email FROM clientes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cliente Cliente
		if err := rows.Scan(&cliente.ID, &cliente.Nome, &cliente.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientes = append(clientes, cliente)
	}
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

// Cria um novo cliente
func createCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//insere no banco de dados
	err := db.QueryRow("INSERT INTO clientes(nome, email) VALUES($1, $2) RETURNING id", cliente.Nome, cliente.Email).Scan(&cliente.ID)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

// Obter pelo ID
func getCliente(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var cliente Cliente
	row := db.QueryRow("SELECT id, nome, email FROM clientes WHERE id = $1", id)
	if err := row.Scan(&cliente.ID, &cliente.Nome, &cliente.Email); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(cliente)
}

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Estrutura do Cliente
type Cliente struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Variável global para o banco de dados
var db *sql.DB

// Função principal
func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Conectar ao banco de dados
	connStr := buildConnectionString()
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verificar a conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Configurar o roteador
	r := mux.NewRouter()

	// Rotas
	r.HandleFunc("/clientes", getClientes).Methods("GET")
	r.HandleFunc("/clientes", createCliente).Methods("POST")
	r.HandleFunc("/clientes/{id}", getCliente).Methods("GET")

	// Iniciar o servidor
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Construir a string de conexão com o banco de dados
func buildConnectionString() string {
	return "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")
}

// Listar todos os clientes
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

// Criar um novo cliente
func createCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Inserir no banco de dados
	err := db.QueryRow("INSERT INTO clientes (nome, email) VALUES ($1, $2) RETURNING id", cliente.Nome, cliente.Email).Scan(&cliente.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

// Obter um cliente pelo ID
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
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

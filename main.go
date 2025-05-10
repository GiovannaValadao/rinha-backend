package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/GiovannaValadao/rinha-backend/db"
)

func main() {
	dbConn := db.Connect()
	defer dbConn.Close()

	r := chi.NewRouter()
	r.Post("/clientes/{id}/transacoes", PostTransacao(dbConn))
	r.Get("/clientes/{id}/extrato", GetExtrato(dbConn))

	log.Println("Rodando na porta 8080")
	http.ListenAndServe(":8080", r)
}

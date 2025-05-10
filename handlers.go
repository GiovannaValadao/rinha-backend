package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func PostTransacao(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var limite, saldo int
		err = db.QueryRow("SELECT limite, saldo FROM clientes WHERE id=$1", id).Scan(&limite, &saldo)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, "internal", 500)
			return
		}

		var t Transacao
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &t); err != nil {
			http.Error(w, "invalid json", 422)
			return
		}
		if t.Valor <= 0 || (t.Tipo != "c" && t.Tipo != "d") || len(t.Descricao) < 1 || len(t.Descricao) > 10 {
			http.Error(w, "invalid fields", 422)
			return
		}

		novoSaldo := saldo
		if t.Tipo == "c" {
			novoSaldo += t.Valor
		} else {
			novoSaldo -= t.Valor
			if novoSaldo < -limite {
				http.Error(w, "limite excedido", 422)
				return
			}
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "internal", 500)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec("UPDATE clientes SET saldo=$1 WHERE id=$2", novoSaldo, id)
		if err != nil {
			http.Error(w, "internal", 500)
			return
		}
		_, err = tx.Exec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1,$2,$3,$4)",
			id, t.Valor, t.Tipo, t.Descricao)
		if err != nil {
			http.Error(w, "internal", 500)
			return
		}
		tx.Commit()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{
			"limite": limite,
			"saldo":  novoSaldo,
		})
	}
}

func GetExtrato(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		var limite, saldo int
		err = db.QueryRow("SELECT limite, saldo FROM clientes WHERE id=$1", id).Scan(&limite, &saldo)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, "internal", 500)
			return
		}

		rows, err := db.Query("SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", id)
		if err != nil {
			http.Error(w, "internal", 500)
			return
		}
		defer rows.Close()
		ultimas := []Transacao{}
		for rows.Next() {
			var t Transacao
			var realizada time.Time
			rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &realizada)
			t.RealizadaEm = realizada.Format(time.RFC3339Nano)
			ultimas = append(ultimas, t)
		}
		dataExtrato := time.Now().Format(time.RFC3339Nano)
		resp := Extrato{}
		resp.Saldo.Total = saldo
		resp.Saldo.DataExtrato = dataExtrato
		resp.Saldo.Limite = limite
		resp.UltimasTransacoes = ultimas

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

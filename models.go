package main

type Cliente struct {
	ID     int
	Limite int
	Saldo  int
}
type Transacao struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}
type Extrato struct {
	Saldo struct {
		Total       int    `json:"total"`
		DataExtrato string `json:"data_extrato"`
		Limite      int    `json:"limite"`
	} `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
}

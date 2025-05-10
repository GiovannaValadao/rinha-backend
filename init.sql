CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100),
    limite INT NOT NULL,
    saldo INT NOT NULL
);

CREATE TABLE IF NOT EXISTS transacoes (
    id SERIAL PRIMARY KEY,
    cliente_id INT NOT NULL REFERENCES clientes(id),
    valor INT NOT NULL,
    tipo CHAR(1) NOT NULL,
    descricao VARCHAR(10) NOT NULL,
    realizada_em TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO clientes (id, nome, limite, saldo) VALUES
  (1, 'o barato sai caro', 100000, 0),
  (2, 'zan corp ltda', 80000, 0),
  (3, 'les cruders', 1000000, 0),
  (4, 'padaria joia', 10000000, 0),
  (5, 'kid mais', 500000, 0);

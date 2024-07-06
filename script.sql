CREATE UNLOGGED TABLE cliente (
  id serial PRIMARY KEY,
  limite integer not null,
  saldo integer not null default 0
);

CREATE UNLOGGED TABLE historico (
  id serial PRIMARY KEY,
  valor integer not null,
  tipo char(1) not null,
  descricao varchar(10) not null,
  client_id integer not null,
  created_at timestamp default(now())
);

CREATE INDEX idx_historico_cliente_id ON historico(client_id);


INSERT INTO cliente (limite, saldo) VALUES
(100000, 0),
(80000, 0),
(1000000, 0),
(10000000, 0),
(500000, 0);

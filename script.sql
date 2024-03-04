DROP TABLE IF EXISTS "transacoes";
DROP TABLE IF EXISTS "clientes";

CREATE UNLOGGED TABLE "clientes" (
    id SERIAL NOT NULL,
    nome TEXT NOT NULL,
    limite INTEGER NOT NULL,
    saldo INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "clientes_pk" PRIMARY KEY ("id")
);

CREATE UNLOGGED TABLE "transacoes" (
    id SERIAL NOT NULL,
    valor INTEGER NOT NULL,
    tipo CHAR(1) NOT NULL,
    descricao VARCHAR(10) NOT NULL,
    realizada_em TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id INTEGER NOT NULL,

    
    CONSTRAINT "transacoes_pk" PRIMARY KEY ("id"),
    CONSTRAINT "transacoes_clientesId_fk" FOREIGN KEY ("client_id") REFERENCES "clientes"("id") ON DELETE RESTRICT ON UPDATE CASCADE
   
);

DO $$
BEGIN
    INSERT INTO clientes (nome, limite)
    VALUES
        ('se voce', 1000 * 100),
        ('ta lendo isso', 800 * 100),
        ('é pq tem', 10000 * 100),
        ('interesse, se tem', 100000 * 100),
        ('ME DA UM ESTAGIO', 5000 * 100);
END; $$

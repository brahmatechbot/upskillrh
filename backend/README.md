# Backend UP Skill

Backend de desenvolvimento em Go para o projeto UP Skill.

## Stack

- Go 1.22+
- PostgreSQL 16+
- pgx/pgxpool
- API REST inicial

## Configuração

Copie `.env.example` se quiser exportar variáveis manualmente:

```bash
export UPSKILLRH_HTTP_ADDR=127.0.0.1:8092
export UPSKILLRH_DATABASE_URL='postgres://upskillrh:upskillrh@localhost:5432/upskillrh_dev?sslmode=disable'
```

## Banco local

```bash
./scripts/setup_dev_db.sh
```

Esse script cria:

- database: `upskillrh_dev`
- user: `upskillrh`
- senha dev: `upskillrh`

## Rodar

```bash
go run .
```

## Testar

```bash
go test ./...
go build ./...
curl http://127.0.0.1:8092/api/health
```

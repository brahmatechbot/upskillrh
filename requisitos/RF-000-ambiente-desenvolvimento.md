# RF-000 - Ambiente de desenvolvimento inicial

## Status
aprovado

## Contexto

O projeto Upskills RH precisa de uma base técnica para evoluir de landing page para aplicação web com backend, banco de dados e frontend estático progressivo.

## Objetivo

Criar um ambiente de desenvolvimento com:

- backend em Go;
- banco PostgreSQL;
- frontend em JavaScript vanilla, CSS e jQuery;
- requisitos documentados em Markdown e gerenciados via Telegram.

## Regras

- O backend fica na pasta `backend`.
- O frontend fica na pasta `site`.
- Os requisitos ficam na pasta `requisitos`.
- O banco local de desenvolvimento se chama `upskillrh_dev`.
- O usuário local de desenvolvimento se chama `upskillrh`.
- A API inicial deve expor `GET /api/health`.
- A porta padrão local do backend deve ser `127.0.0.1:8092` para evitar conflito com outros serviços do servidor.

## Critérios de aceite

- `go test ./...` passa dentro de `backend`.
- `go build ./...` passa dentro de `backend`.
- PostgreSQL local responde ao backend.
- `GET /api/health` retorna JSON com status e estado do banco.
- O site carrega jQuery 4.0.0.

## Histórico

- Criado por Sancho via Telegram em 2026-08-22.

# REQ-20260822 - Tela de login inicial

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 22:11 +0000
- Assunto: requisitosistema-upskillrh tela de login inicial
- Message ID: 1a02b8703c5c4ca4 / Thread ID: 1a02b8703c5c4ca4

## Contexto

Solicitação para implementar a primeira entrega de autenticação do sistema Upskills RH, limitada à tela AUTH-01 — Login, sem alterar a home e sem criar cadastro, recuperação de senha, MFA, SSO, seleção de empresa, logout ou outros módulos.

## Solicitação

Implementar:

- `GET /login` público com tela de login responsiva.
- `POST /api/v1/auth/login` com contrato JSON.
- Template HTML, CSS responsivo e JavaScript vanilla com jQuery.
- Validações de e-mail, senha e `remember_me`.
- Estados de carregamento, erro, credenciais inválidas, rate limit e sucesso.
- Tabelas mínimas `users`, `user_password_credentials` e `auth_sessions`.
- Verificação de senha com Argon2id.
- Sessão opaca em cookie HttpOnly, sem token no JSON.
- CSRF e limite de tentativas.
- Testes para sucesso, senha inválida/rate limit e comportamento técnico relevante.

## Regras e restrições

- Não alterar a home.
- Não implementar cadastro, recuperação de senha, MFA, SSO, seleção de empresa, logout ou outro módulo.
- Não usar mock, usuário fixo, `localStorage` para sessão ou token no JSON.
- Aplicar `trim` somente no e-mail.
- Não alterar/remover espaços da senha.
- Não salvar senha/e-mail em `localStorage`.
- Permitir colar senha.
- Usar `autocomplete="username"` e `autocomplete="current-password"`.
- Não permitir enumeração de usuários em mensagens de credenciais inválidas.
- Não registrar senha ou token em logs.

## Critérios de aceite

- `/login` abre em desktop e mobile.
- Navegação por teclado, labels e mensagens acessíveis.
- Validação de campos no cliente e no servidor.
- Autenticação de usuário real no PostgreSQL.
- Argon2id para verificar senha.
- Sessão opaca em cookie seguro/HttpOnly.
- Prevenção de duplo envio.
- Limite de tentativas.
- Erro genérico para credenciais inválidas.
- Redirecionamento usando `next_url` retornado pelo backend.
- Testes automatizados e build passando.

## Impacto técnico

- Backend Go ganhou módulo `internal/platform/auth` com handler, serviço, repositório Postgres, Argon2id e rate limit em memória.
- Rotas adicionadas ao servidor HTTP: `/login`, `/api/v1/auth/login` e assets `/static/`.
- Adicionadas migrations SQL para tabelas de autenticação.
- Adicionados template, CSS e JS da tela de login sob `backend/web/`.
- Dependência direta adicionada: `golang.org/x/crypto` para Argon2id.

## Ações realizadas

- Criado `backend/web/templates/auth/login.html`.
- Criado `backend/web/static/css/pages/login.css`.
- Criado `backend/web/static/js/pages/login.js`.
- Criados arquivos em `backend/internal/platform/auth/`:
  - `login_handler.go`
  - `login_service.go`
  - `password.go`
  - `postgres_repository.go`
  - `session.go`
  - `login_service_test.go`
- Criadas migrations:
  - `backend/db/migrations/0002_create_auth_tables.up.sql`
  - `backend/db/migrations/0002_create_auth_tables.down.sql`
- Atualizados `backend/app.go`, `backend/main.go`, `backend/go.mod` e `backend/go.sum`.
- Verificação executada:
  - `make test` — passou.
  - `make build` — passou.
  - `curl http://127.0.0.1:8093/login` com backend local — retornou `HTTP/1.1 200 OK` e confirmou conteúdo principal da tela.

## Histórico

- 2026-08-22: requisito recebido por email autorizado, registrado, implementado, verificado localmente e preparado para validação.

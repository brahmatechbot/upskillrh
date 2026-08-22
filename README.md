# upskillrh.com.br

Site institucional estático do projeto Upskills RH.

## Estrutura

- `site/`: conteúdo publicado pelo Nginx para `upskillrh.com.br`.
- `_source_docs/`: materiais extraídos da pasta Google Drive do projeto para embasar copy e posicionamento.

## Fontes analisadas

- `Pitch`
- `Jornada`
- `PesquisaFocadaTecnologia.pdf`

## Posicionamento inicial

Upskills RH ajuda consultorias, integradoras e empresas de tecnologia a reduzir o tempo entre contratação e geração de valor, preparando talentos com diagnóstico, microlearning, validação por evidências, UP DAY e acompanhamento de ramp-up.

## Deploy local

O webroot `/var/www/upskillrh.com.br` aponta para `site/` neste repositório.

## Ambiente de desenvolvimento

### Backend

O backend fica em `backend/` e usa Go com PostgreSQL.

```bash
make db-setup
make test
make build
make run
```

API local padrão:

```text
http://127.0.0.1:8092/api/health
```

### Frontend

O frontend fica em `site/` e usa HTML, CSS, JavaScript vanilla e jQuery 4.0.0 via CDN oficial.

### Requisitos

Os requisitos do produto ficam em `requisitos/` e serão atualizados por solicitações via Telegram.

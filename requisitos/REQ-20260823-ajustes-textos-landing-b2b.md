# REQ-20260823 - Ajustes de textos da landing com foco B2B

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-23 15:58 UTC
- Assunto: requisitosistema-upskillrh ajustes textos landing foco B2B
- Message ID: Gmail message `1a02f58722a0d83b`, thread `1a02f58722a0d83b`

## Contexto

A landing page da UP Skill precisava ser ajustada para uma narrativa mais B2B, conectando seleção, desenvolvimento e onboarding orientados por evidências. A mensagem posiciona a solução para empresas que contratam e desenvolvem em escala, com foco em redução de turnover precoce, aceleração de produtividade e melhor retorno sobre cada contratação.

## Solicitação

Atualizar os textos e a estrutura da landing page para seguir a ordem e o conteúdo enviados por email:

1. Menu: Desafio, Solução, Como funciona, Resultados, Governança, Mercado, Planejar um piloto.
2. Hero com nova promessa: seleção, desenvolvimento e onboarding orientados por evidências.
3. Seção Desafio explicando que o turnover começa antes do desligamento e relacionando gaps, onboarding, sobrecarga dos seniores e IA.
4. Seção Solução: da vaga à autonomia comprovada, com diagnosticar, desenvolver, comprovar e acompanhar.
5. Seção Como funciona com sete etapas entre a vaga e a produtividade, incluindo UP DAY.
6. Seção Resultados com benefícios e indicadores medidos.
7. Seção Governança reforçando IA para escala e decisão humana.
8. Seção Mercado com públicos B2B prioritários.
9. Seção Planejar um piloto e fechamento da página.

## Regras e restrições

- Manter foco B2B e linguagem institucional clara.
- A ordem do menu deve acompanhar exatamente a ordem das seções na página.
- Não executar ações destrutivas, credenciais, cobrança, acessos, banco de dados ou deploy de produção de forma unattended.
- Responder no mesmo thread do Gmail com resumo do que foi feito.

## Critérios de aceite

- Landing page apresenta as seções na ordem solicitada.
- Menu contém os itens informados e aponta para as seções correspondentes.
- Hero, CTAs, seções de conteúdo, indicadores e fechamento refletem o texto recebido.
- Projeto passa em `make test` e `make build`.
- Alterações ficam registradas em commit e enviadas ao GitHub, se possível.

## Impacto técnico

- Alteração principal em `/home/hermes/projects/upskillrh/site/index.html`.
- Sem alteração esperada em backend ou banco de dados.
- Verificação de backend/build via Makefile do projeto e verificação estática do HTML.

## Ações realizadas

- Requisito registrado neste arquivo.
- Landing page atualizada em `site/index.html` com nova narrativa B2B, menu na ordem solicitada, hero, seções Desafio, Solução, Como funciona, Resultados, Governança, Mercado e Planejar um piloto.
- Verificações executadas:
  - `make test` — passou.
  - `make build` — passou.
  - parsing estático do HTML com `html.parser` — passou.
  - verificação local via `python3 -m http.server` e curl com Host header — passou.

## Histórico

- 2026-08-23: Email recebido de Gabriele Guedes e requisito registrado.
- 2026-08-23: Landing atualizada, validada e preparada para commit/push.

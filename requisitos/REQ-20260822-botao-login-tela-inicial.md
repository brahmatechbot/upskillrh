# REQ-20260822 - Botão de login na tela inicial

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 22:25 +0000
- Assunto: requisitosistema-upskillrh incluir botão de login na tela inicial
- Message ID: 1a02b943c3d43969 / Thread ID: 1a02b943c3d43969

## Contexto

A página inicial da UP Skill já possui uma ação principal "Conversar" no cabeçalho e o sistema já tem rota pública de login disponível em `/login`.

## Solicitação

Incluir na tela inicial um botão de login ao lado do botão "Conversar", redirecionando o usuário para a tela de login.

## Regras e restrições

- Manter a ação "Conversar" existente.
- O botão de login deve ficar no cabeçalho, ao lado do botão "Conversar".
- O clique deve redirecionar para `/login`.
- A alteração deve preservar responsividade em telas menores.
- Não envolve mudança destrutiva, credenciais, cobrança, acesso de terceiros ou deploy de produção.

## Critérios de aceite

- A tela inicial exibe um botão "Login" no cabeçalho, ao lado de "Conversar".
- O botão "Login" aponta para `/login`.
- O layout permanece utilizável em desktop e mobile.
- `make test` e `make build` executam com sucesso.

## Impacto técnico

- Alteração de HTML e CSS da landing page em `site/`.
- Reuso da rota de login já existente no backend: `GET /login`.
- Atualização do versionamento do CSS carregado pela página inicial para evitar cache antigo.

## Ações realizadas

- Adicionado agrupamento `nav-actions` no cabeçalho da landing page.
- Incluído botão secundário "Login" com `href="/login"` ao lado do botão "Conversar".
- Ajustados estilos responsivos do cabeçalho para manter os dois botões alinhados.
- Executadas verificações de build/teste do projeto.

## Histórico

- 2026-08-22: requisito recebido por email, registrado, implementado e validado.

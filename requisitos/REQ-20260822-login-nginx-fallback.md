# REQ-20260822 - Correção do acesso ao login pela página inicial

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 22:29 +0000
- Assunto: RE: requisitosistema-upskillrh tela de login inicial
- Message ID: 1a02b977f7fcee30 / Thread ID: 1a02b8703c5c4ca4

## Contexto

Após a inclusão do botão "Login" na tela inicial, o clique atualizava a URL para `/login`, mas a interface permanecia na página inicial. O comportamento indica que a configuração pública do Nginx tratava `/login` como rota estática e aplicava fallback para `/index.html`, em vez de encaminhar a requisição para o backend Go onde a tela de login foi implementada.

## Solicitação

Corrigir o acesso ao botão "Login" da tela inicial para que, ao clicar, a página de login seja realmente renderizada em `/login`.

## Regras e restrições

- Manter o botão "Login" apontando para `/login`.
- Não alterar o escopo funcional do login já implementado.
- Não executar deploy de produção unattended.
- Não alterar credenciais, permissões, banco de dados, cobrança ou serviços de terceiros.

## Critérios de aceite

- `/login` não deve cair no fallback da landing page estática.
- `/login` deve ser encaminhado ao backend Go que renderiza a tela de login.
- `/api/` deve ser encaminhado ao backend para o endpoint de autenticação.
- `/static/` deve servir os assets da tela de login pelo backend.
- `make test` e `make build` devem passar.
- A configuração Nginx deve passar em validação sintática quando aplicável.

## Impacto técnico

- Atualização do arquivo `nginx-upskillrh.com.br.conf` para encaminhar `/login`, `/api/` e `/static/` ao backend local `127.0.0.1:8092` antes do fallback estático da landing page.
- A alteração é de configuração versionada; aplicação em produção ainda depende de atualizar/recarregar o Nginx no servidor.

## Ações realizadas

- Ajustada a configuração Nginx dos blocos HTTP `.com.br` e HTTPS `.online` para proxy reverso de `/login`, `/api/` e `/static/` ao backend.
- Mantido o fallback estático `try_files $uri $uri/ /index.html` apenas para as demais rotas da landing page.
- Executadas verificações do projeto:
  - `make test` — passou.
  - `make build` — passou.
  - `nginx -t` com wrapper/sanitização local da configuração — sintaxe OK e teste bem-sucedido.
  - `curl -H 'Host: upskillrh.com.br' http://127.0.0.1:18080/login` contra Nginx local sanitizado — retornou `HTTP/1.1 200 OK` e renderizou `Entre na sua conta`.
  - `curl -H 'Host: upskillrh.com.br' http://127.0.0.1:18080/static/css/pages/login.css` — retornou `HTTP/1.1 200 OK` e serviu o CSS da tela de login.

## Histórico

- 2026-08-22: bug recebido por email autorizado, registrado, corrigido em configuração versionada e validado tecnicamente.

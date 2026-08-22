# REQ-20260822 - Correção do acesso ao login pela página inicial

## Status
em discussão

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 22:29 +0000; atualização: 2026-08-22 22:37 +0000; nova resposta: 2026-08-22 22:46 +0000
- Assunto: RE: requisitosistema-upskillrh tela de login inicial; RE: RE: requisitosistema-upskillrh tela de login inicial; RE: RE: RE: requisitosistema-upskillrh tela de login inicial
- Message ID: 1a02b977f7fcee30 / Thread ID: 1a02b8703c5c4ca4; atualização Message ID: 1a02b9f13b48c811 / Thread ID: 1a02b8703c5c4ca4; nova resposta Message ID: 1a02ba72f332bc87 / Thread ID: 1a02b8703c5c4ca4

## Contexto

Após a inclusão do botão "Login" na tela inicial, o clique atualizava a URL para `/login`, mas a interface permanecia na página inicial. O comportamento indica que a configuração pública do Nginx tratava `/login` como rota estática e aplicava fallback para `/index.html`, em vez de encaminhar a requisição para o backend Go onde a tela de login foi implementada.

## Solicitação

Corrigir o acesso ao botão "Login" da tela inicial para que, ao clicar, a página de login seja realmente renderizada em `/login`.

Atualização recebida em 2026-08-22 22:37 +0000: a solicitante informou que "o problema continua acontecendo" após a correção versionada. O entendimento técnico é que a alteração já está no repositório, mas ainda precisa ser aplicada/recarregada no Nginx do ambiente público/produção, ação que não deve ser executada unattended sem confirmação/autorização operacional.

Nova resposta recebida em 2026-08-22 22:46 +0000: Gabriele respondeu "Pode autorizar". Como a próxima ação continua sendo operacional em produção (aplicar/recarregar Nginx no ambiente público), a execução unattended permanece bloqueada por política do job; é necessária confirmação/execução explícita de Jeff ou do operador responsável pelo servidor.

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
- Em 2026-08-22 22:37 +0000, retorno da solicitante indicou que o problema ainda aparece no ambiente acessado por ela.
- Em 2026-08-22 22:46 +0000, a solicitante respondeu "Pode autorizar".
- Não foram feitas novas alterações de código nesta rodada; a próxima ação necessária continua sendo confirmar/aplicar a configuração Nginx versionada no ambiente público e recarregar o serviço, pois isso caracteriza ação operacional de produção e permanece bloqueado para execução unattended neste job.

## Histórico

- 2026-08-22: bug recebido por email autorizado, registrado, corrigido em configuração versionada e validado tecnicamente.
- 2026-08-22: novo retorno por email autorizado informou persistência no ambiente; requisito reaberto como "em discussão" e próxima ação registrada como confirmação/aplicação/reload do Nginx em produção, sem execução unattended.
- 2026-08-22: nova resposta autorizando foi registrada; execução em produção segue pendente de Jeff/operador responsável por ser ação operacional fora do escopo unattended do job.

# REQ-20260823 - Ajustes textos landing page

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-23 13:24 UTC
- Assunto: requisitosistema-upskillrh ajustes textos landing page
- Message ID: 1a02ecbb70babe21 / thread 1a02ecbb70babe21

## Contexto

A landing page atual da UP Skill comunica principalmente redução de risco de contratação, bench e retrabalho. A solicitação propõe reposicionar a narrativa para o impacto da IA no primeiro degrau da carreira, destacando a reconstrução de experiências profissionais em ambiente seguro, mensurável e conectado a vagas reais.

## Solicitação

Reformular os textos da home/landing page para cobrir:

- metadados e SEO com o novo posicionamento;
- navegação com Problema, Solução, Jornada, Evidências e Mercado;
- hero sobre a IA ter automatizado o primeiro degrau da carreira;
- problema com quatro efeitos estruturais;
- nova seção de evidências com dados de mercado;
- solução como experiência profissional reconstruída, não apenas capacitação;
- jornada de sete etapas, da vaga ao acompanhamento de 90 dias;
- nova seção “Valor para todos” para candidatos não aprovados e empresas;
- nova seção de governança e uso de IA;
- mercado inicial em consultorias de tecnologia;
- indicadores medidos pela plataforma;
- CTA final de piloto mensurável;
- rodapé com a nova assinatura.

## Regras e restrições

- Manter o escopo na landing page estática do projeto, sem ações destrutivas, produção, credenciais, cobrança, banco de dados ou acesso de terceiros.
- Usar português do Brasil.
- Preservar CTAs existentes de contato por e-mail e login.
- Não inventar links ou fontes além dos nomes fornecidos no email.

## Critérios de aceite

- A landing page deve refletir os novos textos e seções estratégicas solicitadas.
- As seções principais devem ser navegáveis por âncoras.
- Os metadados de SEO/OG devem ser atualizados.
- Deve haver registro rastreável do requisito em Markdown.
- Como houve mudança de frontend estático, o projeto deve ser verificado com `make test`, `make build` e checagem local/renderizada prática.
- As alterações devem ser commitadas e enviadas ao GitHub caso a verificação passe.

## Impacto técnico

- Alteração principal em `/site/index.html`.
- Ajustes de estilos em `/site/styles.css` para suportar novas seções, grids e cards.
- Sem alterações backend esperadas.

## Ações realizadas

- Requisito registrado neste arquivo.
- Landing page atualizada em `/site/index.html` com metadados, navegação, hero, problema, evidências, solução, jornada, valor para todos, governança de IA, mercado, indicadores, CTA final e rodapé.
- Estilos atualizados em `/site/styles.css` para suportar os novos blocos, cards, tabelas e responsividade.
- Logo anexado por Gabriele aplicado na landing page como imagem de marca e favicon: `/site/assets/upskill-logo-gabriele-20260823.png`.
- Logo anexado aplicado inicialmente na tela `/login` como imagem de marca e favicon: `/backend/web/static/img/upskill-logo-gabriele-20260823.png`.
- Ajustes de tamanho/cache bust em `/site/styles.css`, `/site/index.html`, `/backend/web/templates/auth/login.html` e `/backend/web/static/css/pages/login.css`.
- Correção complementar solicitada por Gabriele em 2026-08-23 14:23 UTC: aplicar o mesmo logo também nas demais páginas do fluxo de login/autenticação (`/backend/web/templates/auth/register.html`, `/backend/web/templates/auth/app.html` e `/backend/web/templates/auth/candidate.html`) e uniformizar cache bust do CSS de autenticação.
- Ajuste complementar solicitado por Gabriele em 2026-08-23 14:25 UTC: centralizar os textos do card da página de login/autenticação conforme imagem inline anexada ao email (`image.png`, Content-ID `9dcc0489-e7fb-4d85-b267-5b6c46d7edb0`), preservando campos, checkbox, resumo de erro e formulários alinhados à esquerda para legibilidade.
- Nova observação de Gabriele em 2026-08-23 14:28 UTC: a tela pública de login ainda exibia o logo antigo. Verificação externa em `https://upskillrh.online/login` confirmou que o ambiente publicado ainda estava servindo HTML antigo com `upskills-logo.svg`, enquanto o código no repositório já referencia o PNG novo. Próxima ação depende de autorização explícita de Jeff para restart/deploy/reload em produção; não executado neste job por ser ação de produção vedada sem confirmação.
- Nova verificação da correção complementar executada:
  - `make test` — passou.
  - `make build` — passou.
  - Parser HTML local com `html.parser` — passou para login, cadastro, área da empresa e área do candidato.
  - Checagem estática local — passou: todos os templates de autenticação referenciam `/static/img/upskill-logo-gabriele-20260823.png`, o asset existe e não há referência remanescente a `upskills-logo.svg` nesses templates.
  - Checagem estática local do CSS — passou: `.login-card` usa `text-align: center`, `.brand` centraliza com flex e `.field`, `.check-row` e `.summary` permanecem com `text-align: left`.
  - Verificação pública sem deploy: `https://upskillrh.online/login` ainda retornou HTML antigo com `upskills-logo.svg`; `/static/img/upskill-logo-gabriele-20260823.png` retornou HTTP 200. Diagnóstico: asset novo está disponível, mas a aplicação publicada precisa ser reiniciada/publicada para servir o template atualizado.
- Verificação executada:
  - `make test` — passou.
  - `make build` — passou.
  - Parser HTML local com `html.parser` — passou.
  - Servidor estático local (`python3 -m http.server`) + `curl` com Host header — passou para landing e arquivo PNG do logo.
- Commit enviado ao GitHub: `ca83da2` (`Update landing page positioning copy`).
- Commit enviado ao GitHub para a atualização do logo: `b115a4d` (`Apply updated UP Skill logo`).

## Histórico

- 2026-08-23: solicitação recebida por email de remetente autorizado e registrada para implementação.
- 2026-08-23: alteração de frontend implementada e validada localmente.
- 2026-08-23 14:18 UTC: Gabriele enviou atualização no mesmo thread solicitando considerar a imagem anexada como logo e aplicá-la na página inicial e na página de login. Gmail message ID `1a02efc6da89ba54`, thread `1a02ecbb70babe21`.
- 2026-08-23 14:23 UTC: Gabriele apontou que ainda faltava alterar a página de login. A correção foi interpretada como cobertura completa do fluxo de autenticação, pois `/login` já estava atualizado; registradas alterações em cadastro, área da empresa e área do candidato. Gmail message ID `1a02f00fb4e4b559`, thread `1a02ecbb70babe21`.
- 2026-08-23 14:25 UTC: Gabriele solicitou centralizar os textos da página de login e anexou imagem inline de referência. Gmail message ID `1a02f035ffe25e2d`, thread `1a02ecbb70babe21`.
- 2026-08-23 14:28 UTC: Gabriele reportou que o logo ainda estava errado na tela de login. Foi confirmado que o ambiente público ainda serve template antigo, apesar do repositório já estar corrigido; restart/deploy de produção ficou pendente de autorização explícita de Jeff. Gmail message ID `1a02f05f744f488b`, thread `1a02ecbb70babe21`.

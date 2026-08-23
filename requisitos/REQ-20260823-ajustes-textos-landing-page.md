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
- Verificação executada:
  - `make test` — passou.
  - `make build` — passou.
  - Parser HTML local com `html.parser` — passou.
  - Servidor estático local (`python3 -m http.server`) + `curl` com Host header — passou.
- Commit enviado ao GitHub: `ca83da2` (`Update landing page positioning copy`).

## Histórico

- 2026-08-23: solicitação recebida por email de remetente autorizado e registrada para implementação.
- 2026-08-23: alteração de frontend implementada e validada localmente.

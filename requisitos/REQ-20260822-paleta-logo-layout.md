# REQ-20260822 - Paleta, logo e layout conforme imagem de referência

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 23:03 UTC; atualização em 2026-08-22 23:16 UTC; atualização em 2026-08-22 23:23 UTC; atualização em 2026-08-22 23:28 UTC
- Assunto: requisitosistema-upskillrh redefinir padrao de cores e logo; Re: requisitosistema-upskillrh redefinir padrao de cores e logo
- Message ID: 1a02bb6b4a516592 / thread 1a02bb6b4a516592; atualização 1a02bc295fbf3a7d / thread 1a02bb6b4a516592; atualização 1a02bc9308a83d5a / thread 1a02bb6b4a516592; atualização 1a02bcdf588ebba1 / thread 1a02bb6b4a516592

## Contexto

A solicitação inicial pediu que a imagem enviada por email fosse usada como referência para redefinir a paleta de cores, o logo e o layout visual da landing page da UP Skill.

Na atualização, a remetente informou que o resultado ainda não ficou adequado e pediu aplicação do padrão de cores em todas as páginas, incluindo a tela de login, além de recriação das imagens respeitando cada contexto e a nova paleta. Em nova resposta, apontou especificamente que ainda faltava alterar a imagem da seção "Foco inicial: consultorias de tecnologia". Na atualização mais recente, corrigiu a nomenclatura oficial da plataforma: o nome correto é "UP Skill", não "UP Skills RH", "Upskills RH" ou "UPSKILLING".

## Solicitação

- Considerar a imagem anexada como fonte para definição de paleta de cores e logo.
- Rever todo o layout público para seguir o padrão visual da imagem enviada.
- Aplicar o padrão de cores também na tela de login.
- Recriar as imagens/ilustrações do site respeitando o contexto de cada seção e a paleta preta, lavanda, roxa e amarelo-limão.
- Alterar especificamente a imagem da seção "Foco inicial: consultorias de tecnologia" para uma ilustração contextual aderente à nova identidade.
- Corrigir a nomenclatura pública e institucional da plataforma para "UP Skill".
- Usar o anexo `Image.jpeg` como referência e asset local versionado.

## Regras e restrições

- Escopo limitado ao frontend público em `/site` e à tela de login do backend em `/backend/web/templates/auth/login.html` e `/backend/web/static/css/pages/login.css`.
- Não executar deploy de produção, alterações de acesso, credenciais, banco de dados, cobrança ou integrações externas de forma autônoma.
- Manter conteúdo, navegação, acessibilidade básica e responsividade existentes.
- Versionar os assets visuais para rastreabilidade.

## Critérios de aceite

- A página pública usa paleta coerente com a referência: base escura/preta, lavanda, roxo e amarelo-limão.
- A tela de login usa a mesma identidade visual.
- O logo SVG é atualizado para refletir a nova identidade visual.
- A imagem hero é recriada como ilustração contextual da jornada UP Skill.
- O mapa visual de mercado é recriado/ajustado para a mesma paleta.
- A imagem da seção "Foco inicial: consultorias de tecnologia" deixa de ser genérica e comunica consultorias, talentos, clientes, squads, evidências, ramp-up e margem.
- A página pública, tela de login, logos, imagens SVG, README e registros de requisitos não exibem mais a plataforma como "UP Skills RH", "Upskills RH" ou "UPSKILLING".
- Cache-busting dos assets estáticos é atualizado.
- `make test` e `make build` passam.
- Verificação estática/local confirma que HTML, CSS e assets principais são servidos corretamente.

## Impacto técnico

- Frontend estático em `/site`.
- Tela de login renderizada pelo backend Go.
- Arquivos alterados/criados:
  - `site/index.html`
  - `site/styles.css`
  - `site/upskills-logo.svg`
  - `site/market-map.svg`
  - `site/assets/brand-reference-gabriele-20260822.jpeg`
  - `site/assets/hero-upskills-context.svg`
  - `backend/web/templates/auth/login.html`
  - `backend/web/static/css/pages/login.css`
  - `backend/web/static/img/upskills-logo.svg`

## Ações realizadas

- Baixei o anexo do Gmail e salvei em `site/assets/brand-reference-gabriele-20260822.jpeg`.
- Atualizei a referência da imagem principal da home para o novo asset na primeira entrega.
- Redefini variáveis e aplicações de cor no CSS da home para uma identidade visual baseada em preto, lavanda, roxo e amarelo-limão.
- Atualizei o logo SVG para usar preto, roxo, lavanda e amarelo-limão.
- Atualizei versões de CSS/JS no HTML para reduzir risco de cache antigo.
- Na atualização `1a02bc295fbf3a7d`, apliquei a mesma paleta à tela de login.
- Substituí o marcador textual da marca no login pelo logo SVG oficial.
- Recriei a imagem hero como SVG contextual (`site/assets/hero-upskills-context.svg`) com diagnóstico, microlearning, UP DAY e evidências.
- Recriei/ajustei o mapa de mercado (`site/market-map.svg`) para a mesma identidade visual.
- Na atualização `1a02bc9308a83d5a`, refiz a imagem da seção "Foco inicial: consultorias de tecnologia" como uma ilustração contextual, conectando talentos, clientes, squads, evidências, ramp-up, bench e margem.
- Atualizei a referência da imagem no HTML para `market-map.svg?v=2` e melhorei o texto alternativo para refletir o novo contexto.
- Na atualização `1a02bcdf588ebba1`, padronizei o nome público da plataforma como "UP Skill" em títulos, metas, CTAs, textos, alt text, SVGs, tela de login, README e documentação interna relevante.
- Mantive os nomes físicos de assets já versionados (`upskills-logo.svg` e `hero-upskills-context.svg`) para evitar quebra de referências, alterando apenas o conteúdo/rotulagem visível.

## Histórico

- 2026-08-22 23:03 UTC: requisito recebido por email autorizado, registrado e implementado na landing page.
- 2026-08-22 23:16 UTC: ajuste recebido por email autorizado; requisito atualizado para ampliar padrão visual à tela de login e recriar imagens contextuais.
- 2026-08-22 23:23 UTC: ajuste recebido por email autorizado; imagem da seção "Foco inicial: consultorias de tecnologia" recriada como ilustração contextual e referência atualizada no HTML.
- 2026-08-22 23:28 UTC: ajuste recebido por email autorizado; nomenclatura pública e institucional padronizada para "UP Skill" e verificada no site/login.

# REQ-20260822 - Paleta, logo e layout conforme imagem de referência

## Status
validado

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-22 23:03 UTC
- Assunto: requisitosistema-upskillrh redefinir padrao de cores e logo
- Message ID: 1a02bb6b4a516592 / thread 1a02bb6b4a516592

## Contexto

A solicitação pede que a imagem enviada por email seja usada como referência para redefinir a paleta de cores, o logo e o layout visual da landing page da Upskills RH.

## Solicitação

- Considerar a imagem anexada como fonte para definição de paleta de cores e logo.
- Rever todo o layout público para seguir o padrão visual da imagem enviada.
- Usar o anexo `Image.jpeg` como referência e asset local versionado.

## Regras e restrições

- Escopo limitado ao frontend público em `/site`.
- Não executar deploy de produção, alterações de acesso, credenciais, banco de dados, cobrança ou integrações externas de forma autônoma.
- Manter conteúdo, navegação e responsividade existentes.
- Versionar o asset recebido para rastreabilidade.

## Critérios de aceite

- A página passa a usar paleta coerente com a imagem: base escura/preta, lavanda, roxo e amarelo-limão.
- O logo SVG é atualizado para refletir a nova identidade visual.
- A imagem hero passa a usar o novo anexo recebido.
- Cache-busting dos assets estáticos é atualizado.
- `make test` e `make build` passam.
- Verificação estática/local confirma que HTML, CSS e assets principais são servidos corretamente.

## Impacto técnico

- Frontend estático em `/site`.
- Arquivos alterados/criados:
  - `site/index.html`
  - `site/styles.css`
  - `site/upskills-logo.svg`
  - `site/assets/brand-reference-gabriele-20260822.jpeg`

## Ações realizadas

- Baixei o anexo do Gmail e salvei em `site/assets/brand-reference-gabriele-20260822.jpeg`.
- Atualizei a referência da imagem principal da home para o novo asset.
- Redefini variáveis e aplicações de cor no CSS para uma identidade visual baseada em preto, lavanda, roxo e amarelo-limão.
- Atualizei o logo SVG para usar preto, roxo e amarelo-limão.
- Atualizei versões de CSS/JS no HTML para reduzir risco de cache antigo.

## Histórico

- 2026-08-22: requisito recebido por email autorizado, registrado, implementado e preparado para validação técnica.

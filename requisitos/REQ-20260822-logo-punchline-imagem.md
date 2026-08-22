# REQ-20260822 - Logo, punch line e imagem principal

## Status
validado

## Origem

- Canal: email
- Remetente: Isabela Ferraz <isabelaferraz.1997@gmail.com>
- Data do email: 2026-08-22 19:46 -0300
- Assunto: requisitosistema-upskillrh - Logo e punch line
- Message ID: 1a02babcf614409a / thread 1a02babcf614409a

## Contexto

A solicitação pede atualização da comunicação visual inicial da landing page da Upskills RH, trocando a punch line atual e usando a imagem enviada em anexo como referência visual da interface principal.

## Solicitação

- Substituir a frase de destaque `Contratar é caro. Demorar para gerar valor custa mais.` por `POTENCIAL QUE SE PROVA`.
- Mudar a imagem/visual principal da interface para o arquivo anexado ao email: `WhatsApp Image 2026-08-22 at 00.50.00.jpeg`.

## Regras e restrições

- Manter o escopo na landing page pública.
- Não alterar credenciais, acessos, banco de dados, deploy de produção ou integrações externas.
- Preservar responsividade e estrutura visual existente.
- Usar o arquivo recebido como asset local do site, sem depender de link externo do email.

## Critérios de aceite

- A home exibe a nova punch line: `Potencial que se prova.`.
- A imagem principal da seção hero usa o anexo recebido por email.
- O asset está versionado dentro do projeto.
- `make test` e `make build` passam.
- Verificação estática confirma que o HTML referencia o novo asset.

## Impacto técnico

- Frontend estático em `/site`.
- Arquivos alterados:
  - `site/index.html`
  - `site/styles.css`
  - `site/assets/hero-reference-isabela-20260822.jpeg`

## Ações realizadas

- Baixei o anexo do Gmail e salvei em `site/assets/hero-reference-isabela-20260822.jpeg`.
- Atualizei o H1 da seção hero para `Potencial que se prova.`.
- Substituí a imagem hero anterior pelo asset enviado.
- Ajustei a proporção do card de imagem principal para 3:2, compatível com o anexo recebido.

## Histórico

- 2026-08-22: requisito recebido por email autorizado, registrado, implementado e preparado para validação técnica.

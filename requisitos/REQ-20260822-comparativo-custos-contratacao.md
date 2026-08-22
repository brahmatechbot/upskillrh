# REQ-20260822 - Comparativo de custos de contratação

## Status
validado

## Origem

- Canal: email
- Remetente: Marcio Farinasso <marciofarinasso@gmail.com>
- Data do email: 2026-08-22 18:39 -0300
- Assunto: Requisitosistema-upskillrh contratação
- Message ID: 1a02b69f9b6a1a7d / thread 1a02b69f9b6a1a7d

## Contexto

Solicitação enviada por membro autorizado do projeto UP Skill para complementar o site com conteúdo comparativo sobre custo de contratação e atualizar o endereço de contato exibido nos CTAs.

## Solicitação

- Criar uma tabela de comparação de custo de contratação da área de tecnologia e da área da indústria.
- Alterar no site o email `jeffprestes@gmail.com` para `brahmatechbot@gmail.com`.

## Regras e restrições

- Manter a alteração limitada ao site estático do projeto.
- Não executar ações destrutivas, alterações de acesso, credenciais, banco de dados, billing ou deploy de produção sem autorização explícita.
- Manter conteúdo em português e alinhado ao posicionamento da UP Skill.

## Critérios de aceite

- O site apresenta uma tabela clara comparando fatores de custo de contratação entre tecnologia e indústria.
- Todos os links `mailto:` que apontavam para `jeffprestes@gmail.com` passam a apontar para `brahmatechbot@gmail.com`.
- A página continua válida como HTML estático e sem quebra visual evidente.
- `make test` e `make build` executam com sucesso.

## Impacto técnico

- Frontend estático em `site/index.html` e `site/styles.css`.
- Sem alteração no backend, banco de dados ou infraestrutura.

## Ações realizadas

- Criada seção `Custo de contratação: tecnologia x indústria` com tabela comparativa no site.
- Adicionados estilos responsivos para a tabela, com rolagem horizontal em telas menores.
- Atualizados os CTAs de email do site de `jeffprestes@gmail.com` para `brahmatechbot@gmail.com`.
- Verificações executadas: `make test`, `make build` e validação local com servidor HTTP estático/curl.

## Histórico

- 2026-08-22: Requisito recebido por email, classificado como alteração de conteúdo/site segura e implementado.

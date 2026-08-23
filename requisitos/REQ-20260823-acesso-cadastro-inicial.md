# REQ-20260823 - Acesso e cadastro inicial

## Status
em desenvolvimento

## Origem

- Canal: email
- Remetente: Gabriele Guedes <gabrielepgsantos@outlook.com>
- Data do email: 2026-08-23 12:36 +0000
- Assunto: requisitosistema-upskillrh Pacote de desenvolvimento Acesso e cadastro inicial
- Message ID: 1a02e9f3ac742372 / Thread ID: 1a02e9f3ac742372

## Contexto

Solicitação de pacote para permitir que visitantes acessem login, escolham cadastro empresarial ou de candidato, criem conta real no PostgreSQL, sejam autenticados automaticamente e acessem a área logada correspondente sem mocks.

## Solicitação

Implementar o pacote de autenticação e cadastro inicial composto por:

- `AUTH-01 — Login`: tela `/login` com opções de criação de conta empresarial e de candidato, ambas apontando para `/cadastro` com parâmetro de tipo; login deve autenticar contas dos dois contextos e redirecionar conforme perfis/vínculos reais.
- `AUTH-02 — Cadastro mínimo`: página `/cadastro` com alternância entre empresa e candidato, campos mínimos, validação de senha, CPF/CNPJ, termos, privacidade e marketing; persistência transacional de usuário, credencial, perfil/organização/vínculo, aceites, marketing, sessão e token de verificação.
- `CAND-01 — Área inicial do candidato`: página protegida `/candidate` com saudação, confirmação de conta, aviso de e-mail pendente e próximos passos sem links quebrados.
- `COMP-01 — Área inicial da empresa`: página protegida `/app` com saudação, dados básicos da organização, status, próximos passos e indicadores iniciais sem links quebrados.

## Regras e restrições

- Não criar `user_type` mutuamente exclusivo; contexto deve ser resolvido por `candidate_profiles` e `organization_memberships`.
- CPF/CNPJ devem ser normalizados, validados no backend, protegidos em armazenamento e indexados por índice cego para duplicidade.
- Senhas devem usar Argon2id e não devem aparecer em logs.
- Termos, política de privacidade e marketing devem ser registros separados.
- Sessão deve ser criada por cookie seguro/HttpOnly após cadastro.
- Cadastro deve ocorrer em transação única.
- Área de candidato não deve exibir CPF.
- Área empresarial não deve exibir CNPJ integral.
- Não criar telas futuras nem links quebrados para vagas, currículo, experiências ou candidaturas.
- Em ambiente local, verificação de e-mail é registrada como token pendente; não foi criado bypass de produção.

## Critérios de aceite

- Login apresenta duas possibilidades de cadastro.
- Opções direcionam para `/cadastro?tipo=empresa` e `/cadastro?tipo=candidato`.
- Cadastro alterna visualmente entre empresa e candidato e respeita parâmetro inicial.
- Cadastro empresarial persiste usuário, credencial, organização, membership, aceites, marketing, sessão e verificação pendente.
- Cadastro de candidato persiste usuário, credencial, perfil básico, aceites, marketing, sessão e verificação pendente.
- CPF e CNPJ são validados no backend e armazenados sem plaintext.
- Termos, privacidade e marketing são persistidos separadamente.
- Login posterior redireciona candidato para `/candidate` e usuário empresarial para `/app` com base nos vínculos reais.
- `/candidate` e `/app` exigem sessão e carregam contexto pelo token de sessão no servidor.
- Build e testes do backend passam.

## Impacto técnico

- Novas rotas backend: `/cadastro`, `/candidate`, `/app`, `/api/v1/auth/register`.
- Login passa a resolver destino por perfil/vínculo real.
- Novas tabelas/migração para cadastro mínimo: `organizations`, `organization_memberships`, `candidate_profiles`, `industry_segments`, `policy_versions`, `user_legal_acceptances`, `marketing_consents`, `email_verification_tokens` e `users.email_verified_at`.
- Novos templates e JavaScript/CSS para fluxo de cadastro e áreas iniciais.
- Identificadores fiscais são protegidos por índice cego HMAC usando `UPSKILLRH_PII_SECRET` quando configurado, com fallback local para desenvolvimento.

## Ações realizadas

- Registrado este requisito em `requisitos/REQ-20260823-acesso-cadastro-inicial.md`.
- Atualizada tela de login para incluir chamadas de cadastro empresarial e de candidato.
- Criada página `/cadastro` com alternância entre tipos, campos mínimos, aceite legal e marketing opcional.
- Criada API de cadastro com validação backend para CPF/CNPJ, senha, campos específicos e aceites obrigatórios.
- Criada persistência transacional para cadastro de candidato e empresa.
- Criada migração `0003_registration_flow` com tabelas e catálogos necessários.
- Criadas páginas protegidas `/candidate` e `/app` com próximos passos desabilitados/identificados como etapas futuras.
- Ajustado login para redirecionar por vínculo real de candidato/empresa.
- Verificação executada:
  - `make test` — passou.
  - `make build` — passou.

## Histórico

- 2026-08-23: Requisito recebido por email autorizado, registrado e implementado em primeira versão funcional do pacote de cadastro inicial.

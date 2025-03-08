# Sistema de Cadastro de Clientes

Este é um microsserviço desenvolvido em **Golang** para gerenciar o cadastro de clientes. O sistema permite que os usuários se cadastrem fornecendo informações básicas, como nome, e-mail, telefone e endereço.

## Funcionalidades Atuais

- **Cadastro de Clientes:**
  - Os usuários podem se cadastrar fornecendo:
    - Nome
    - Sobrenome
    - E-mail
    - Telefone
    - Endereço (CEP, Cidade, Bairro, Rua, Número, Complemento)
  - Validação de e-mail e CEP usando APIs externas.

## Tecnologias Utilizadas

- **Backend:** Golang
- **Banco de Dados:** PostgreSQL
- **Autenticação:** JWT (JSON Web Tokens)

## Como Executar o Projeto

### Pré-requisitos

- Docker e Docker Compose instalados.
- Go 1.20 ou superior.

### Passos para Execução

1. Clone o repositório:
   ```bash
   git clone https://github.com/seu-usuario/sistema-cadastro-clientes.git
   cd sistema-cadastro-clientes
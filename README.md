# kn

## Descrição

O **kn** é um sistema em desenvolvimento que visa se tornar uma plataforma completa para e-commerce. Atualmente, o projeto possui apenas a funcionalidade de autenticação de produtos pronta, mas futuras atualizações incluirão diversas funcionalidades para gerenciar uma loja online. O projeto está sendo construído com o menor número de pacotes externos possível, focando em uma base de código limpa e de fácil manutenção.

## Instalação

Passos para instalar e configurar o projeto:

1. Clone o repositório:

   ```bash
   git clone https://github.com/seu-usuario/kn-server.git
   ```

2. Entre no diretório do projeto:

   ```bash
   cd kn-server
   ```

3. Copie o arquivo de exemplo `.env.example` para `.env` e preencha com suas configurações:

   ```bash
   cp .env.example .env
   ```

4. Inicie o banco de dados via Container Docker:

   ```bash
   docker-compose up
   ```

## Uso

Instruções básicas de como usar o projeto após a instalação:

1. Para iniciar o projeto:

   ```bash
   make run
   ```

2. Para rodar os testes:

   ```bash
   make test
   ```

## Recursos

Recursos adicionais podem ser adicionados aqui no futuro.

## Autores

- **Rafael Bernardino**

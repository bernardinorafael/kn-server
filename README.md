# kn

## Description

**kn** is a system under development that aims to become a complete platform for e-commerce. Currently, the project has the authentication system, which is receiving updates, future updates will include several features for managing an online store. The project is being built with as few external packages as possible, focusing on a clean and maintainable code base.

## Install

Steps to install and configure the project:

1. Clone repo:

   ```bash
   git clone https://github.com/seu-usuario/kn-server.git
   ```

2. Enter the project directory:

   ```bash
   cd kn-server
   ```

3. Copy the example file `.env.example` to `.env` and fill in your settings:

   ```bash
   cp .env.example .env
   ```

4. Start the database via Docker Container:

   ```bash
   docker-compose up
   ```

## Uso

Basic instructions on how to use the project after installation:

1. To start the project:

   ```bash
   make run
   ```

2. To run the tests:

   ```bash
   make test
   ```

3. To enter the database inside Docker container:
   ```bash
   make pg
   ```

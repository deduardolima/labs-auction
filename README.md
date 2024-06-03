# Labs Auction

## Descrição

Este projeto simula um sistema de leilão utilizando a arquitetura limpa na linguagem Go, com suporte para servidor Web. A funcionalidade de leilão permite criar leilões, fazer lances, e determinar automaticamente o fechamento dos leilões após um tempo definido.

## Funcionalidades Implementadas
- Criação de Leilão: Permite criar novos leilões com informações sobre o produto, categoria, descrição, condição e tempo de duração.
- Lances: Permite que os usuários façam lances em leilões abertos.
- Fechamento Automático de Leilões: Utilizando goroutines, o sistema verifica continuamente os leilões e fecha automaticamente aqueles cujo tempo de duração expirou.

## Configuração
```
git clone https://github.com/deduardolima/labs-auction.git
cd labs-auction

```

## Instalação e Execução com Docker
Construa e inicie os containers:
```
docker-compose up --build -d
```

isso irá construir a imagem do aplicativo e iniciar os serviços definidos no docker-compose.yml, incluindo o banco de dados e o aplicativo.

A aplicação estará acessível em http://localhost:8080

## Criação de Usuário

Após iniciar os containers, é necessário criar um usuário manualmente. Siga os comandos abaixo:

### Acesse o Mongo Express
1. Abra o Mongo Express no seu navegador: [Mongo-Express](http://localhost:8081)
2. Faça login com as credenciais:
- **Username**: admin
- **Password**: admin

### Crie o Banco de Dados e a Coleção
1. Crie o banco de dados `auctions`.
2. Crie uma coleção chamada `users` dentro do banco de dados `auctions`.

### Insira um Documento na Coleção `users`
1. Insira o seguinte documento na coleção `users`:
```json
{
  "_id": "860c93b4-0282-451e-8249-b55ee52d1460",
  "name": "John Doe"
}
```
2. Clique em `Save` para salvar o documento.

## Execucão 

No arquivo `api/auction.http` fornece os endpoints para criação de leilão, lances e conferir quem ganhou leilão com ID auction por exemplo.

Atualmente o tempo limite para fechamento do leilão é de 5 minutos, pode ser alterado no arquivo `.env` na variavel `AUCTION_DURATION`



### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis exemplo :

- BATCH_INSERT_INTERVAL=20s
- MAX_BATCH_SIZE=4
- AUCTION_INTERVAL=20s
- AUCTION_DURATION=5
- MONGO_INITDB_ROOT_USERNAME=admin
- MONGO_INITDB_ROOT_PASSWORD=admin
- MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
- MONGODB_DB=auctions


# Descrição desafio

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos, banco de dados e manipulação de arquivos com Go.
 
Você precisará nos entregar dois sistemas em Go:
- client.go
- server.go
 
Os requisitos para cumprir este desafio são:
 
O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
 
Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
 
O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
 
Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
 
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.


### Configuração do Ambiente

Este projeto utiliza variáveis de ambiente para configurar aspectos importantes da aplicação, como acesso ao banco de dados, chaves de API e outras configurações sensíveis.

#### Passo 1: Copiar o arquivo `env.sample` para `.env`

Antes de rodar o projeto, você precisa criar o arquivo `.env` na raiz do projeto. Para isso, copie o arquivo `env.sample` para `.env`:


#### Variáveis de Ambiente

Este projeto utiliza variáveis de ambiente para configurar as diferentes partes do sistema, como o caminho do banco de dados, a URL do servidor e o local de armazenamento de arquivos. Abaixo está a descrição de cada variável de ambiente que você precisa configurar no arquivo `.env`.

#### Descrição das Variáveis de Ambiente

##### 1. `DATABASE_PATH`

- **Descrição**: Caminho para o arquivo de banco de dados.
- **Valor esperado**: Um caminho relativo ou absoluto onde o banco de dados SQLite será armazenado.
- **Exemplo**:

  ```env
  DATABASE_PATH=../data/app.db
  ```

##### 2. `SERVER_HOST`

- **Descrição**: URL do servidor, incluindo o host e a porta em que o servidor Go está sendo executado.
- **Valor esperado**: URL com o protocolo (http ou https), seguido pelo host e porta.
- **Exemplo**:

  ```env
  SERVER_HOST=http://localhost:8080
  ```

##### 3. `FILE_PATH_STORAGE`

- **Descrição**:  Caminho para o diretório onde os arquivos serão armazenados.
- **Valor esperado**:  Caminho relativo ou absoluto para o diretório onde os arquivos da aplicação serão salvos.
- **Exemplo**:

  ```env
  FILE_PATH_STORAGE=..
  ```

### Solução para "Erro ao Conectar ao Banco de Dados" no Windows

Se você está recebendo o seguinte erro ao tentar rodar o aplicativo no Windows:

> failed to initialize database, got error Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub 2025/03/24 15:28:18 Erro ao conectar no banco de > dados: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub exit status 1


Esse erro ocorre porque o driver `go-sqlite3` requer o **cgo**, que pode estar desabilitado por padrão na sua compilação. Para corrigir esse problema, acesse o link abaixo e siga os passos fornecidos na discussão para resolver o erro:

> [https://kodekloud.com/community/t/issue-while-running-go-webapp-sample-application-in-local-machine/363181/2](https://kodekloud.com/community/t/issue-while-running-go-webapp-sample-application-in-local-machine/363181/2)


## Rodando o Servidor

Para rodar o servidor, siga os passos abaixo.

### Passo 1: Configuração do Ambiente

Antes de rodar o servidor, certifique-se de que o arquivo `.env` está configurado corretamente. Consulte a seção de [Variáveis de Ambiente](#variáveis-de-ambiente) para detalhes sobre como preencher o arquivo `.env`.

### Passo 2: Rodar o Servidor

Com o ambiente configurado e as variáveis de ambiente preenchidas, você pode iniciar o servidor executando o comando abaixo:

```bash
go run server/server.go
```

### Passo 3: Verificar os Logs

Ao rodar o servidor, você verá logs de inicialização no terminal. O servidor vai se conectar ao banco de dados e exibir mensagens como a seguir:

> 
> Inicializando conexão com banco de dados
> 
> Conexão estabelecida
> 


Após o servidor estar rodando, você pode fazer requisições para a API. O servidor estará disponível na URL configurada na variável de ambiente SERVER_HOST.


## Rodando o Cliente

Este projeto inclui um cliente que faz requisições para o servidor e exibe as cotações das moedas. Abaixo estão as instruções para rodar o cliente.


### Passo 1:  Rodar o Cliente

#### 1.1. Navegar até o Diretório client

Primeiro, navegue até o diretório client onde o arquivo client.go está localizado:

```bash
cd client
```
#### 1.2. Rodar o Cliente

Com o ambiente configurado, execute o seguinte comando para rodar o cliente:

```bash
go run client.go
```

#### 1.3.  Verificar os Logs e Resposta

>
> Arquivo criado com sucesso!
>
> Cotação de USD para BRL: 5.7499
>  


## Conferindo a Tabela `exchanges` no Banco de Dados SQLite

Após rodar o cliente e o servidor, você pode verificar as informações armazenadas na tabela `exchanges` dentro do banco de dados SQLite. O banco de dados está localizado no caminho configurado na variável `DATABASE_PATH` no arquivo `.env`.

### Passo 1: Localizar o Banco de Dados

O caminho do banco de dados é configurado pela variável `DATABASE_PATH`. No seu arquivo `.env`, você verá algo como:

```env
DATABASE_PATH=../data/app.db
```

### Passo 2: Verificar os Dados na Tabela exchanges

### Exemplo de Estrutura da Tabela `exchanges`

Após rodar o cliente e o servidor, os dados das cotações são armazenados na tabela `exchanges` do banco de dados SQLite. A seguir, um exemplo de como a tabela `exchanges` pode ser estruturada:

| **ID** | **from_currency** | **to_currency** | **rate** | **retrieval_date**        |
|--------|-------------------|-----------------|----------|---------------------------|
| 1      | USD               | BRL             | 5.7499    | 2025-03-24 12:34:56       |
| 2      | EUR               | USD             | 5.7521    | 2025-03-24 12:35:00       |

#### Descrição dos Campos

- **ID**: Identificador único para cada registro da tabela.
- **from_currency**: A moeda de origem para a cotação.
- **to_currency**: A moeda de destino para a cotação.
- **rate**: A taxa de câmbio entre as moedas de origem e destino.
- **retrieval_date**: A data e hora em que a cotação foi obtida.

#### Como a Tabela é preenchida

A tabela `exchanges` é preenchida automaticamente quando o cliente faz uma requisição ao servidor para obter a cotação das moedas. O cliente envia a requisição para o servidor, que então processa os dados e armazena as cotações obtidas no banco de dados.

Se você verificar a tabela, verá registros semelhantes a este exemplo, com as cotações obtidas e as respectivas datas de consulta.

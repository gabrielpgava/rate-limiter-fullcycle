# Rate Limiter (Go) - Go Expert

A aplicação é um rate limiter que analisa conexão por IP e IP + Token, os bloqueios padrões acontecem quando atinge 10 req/s e com token a 100 req/s

## Variaveis de ambiente

- `API_KEY`: token valido para o header `API_KEY`.
- `STORAGE_PROVIDER`: `redis` (padrao) ou `memory`.
- `REDIS_ADDR`: endereco do Redis (ex: `localhost:6379`).
- `REDIS_PASSWORD`: senha do Redis (opcional).
- `max_request_ip_per_second`: limite por IP (default 5 se nao setado).
- `max_request_token_per_second`: limite por token (default 5 se nao setado).
- `timeout_ip_block_inSeconds`: tempo de bloqueio por IP (segundos).
- `timeout_token_block_inSeconds`: tempo de bloqueio por token (segundos).


## Executar localmente

1. Suba o Redis:
   ```bash
   docker run --rm -p 6379:6379 redis:latest
   ```
2. Rode a aplicacao:
   ```bash
   go run ./cmd
   ```

Servidor sobe em `:8080`.

## Executar com Docker Compose

```bash
docker-compose up --build
```

Teste rapido de carga (usa `ghcr.io/hatoo/oha:latest`):

```bash
docker-compose --profile loadtest up --build
```

## Exemplos de Use

Use o requests.http no visual code para testar os endpoints de forma manual


## Testes

Unitarios:

```bash
go test ./...
```

Integracao (requer Redis em `REDIS_ADDR`):

```bash
go test ./... -tags=integration
```

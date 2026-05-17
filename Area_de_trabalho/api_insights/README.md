# API de Insights

API Go simples para gerar um relatório de insights a partir de mocks de vendas, sugestões, reclamações e dados captados dos usuários.

## Arquitetura

- `cmd/api`: composição da aplicação e servidor HTTP.
- `internal/model`: entidades, repositório mock, normalização e geração do relatório.
- `internal/controller`: rotas HTTP e orquestração das requisições.
- `internal/view`: renderização das respostas JSON.

### Diagrama Arquitetural

```mermaid
graph TD
	subgraph cmd/api
		A[main.go] -->|Cria| B[NewMockStore]
		A -->|Cria| C[NewInsightModel]
		A -->|Cria| D[NewInsightController]
		D -->|Expõe| E[Routes]
	end
	subgraph Model
		B[NewMockStore] --> F[MockStore]
		F -->|Implementa| G[InsightRepository]
		C[NewInsightModel] --> H[InsightModel]
		H -->|Usa| G
		H -->|Gera| I[GenerateReport]
		H -.-> J[Entidades e DTOs]
	end
	subgraph Controller
		D[NewInsightController] --> E[Routes]
		E -->|Define| K[GET /health e GET /insights/report]
		K -->|Chama| I
	end
	subgraph View
		K -->|Renderiza| L[JSON e Error]
	end
```

## Rodando

```bash
go run ./cmd/api
```

Endpoints:

- `GET /health`
- `GET /insights/report`

## Testes

```bash
go test ./...
go test -race ./...
```

# API de Insights

API Go simples para gerar um relatório de insights a partir de mocks de vendas, sugestões, reclamações e dados captados dos usuários.

## Arquitetura

- `cmd/api`: composição da aplicação e servidor HTTP.
- `internal/domain`: modelos de domínio e DTO do relatório.
- `internal/repository`: mock simulando banco de dados, protegido com `sync.RWMutex` e cópias defensivas.
- `internal/service`: normalização, tratamento de dados, geração de insights e uso de goroutines.
- `internal/transport/http`: handlers HTTP.

### Diagrama Arquitetural

```mermaid
graph TD
	subgraph cmd/api
		A[main.go] -->|Cria| B[NewMockStore]
		A -->|Cria| C[NewInsightService]
		A -->|Cria| D[NewHandler]
		D -->|Expõe| E[Routes]
	end
	subgraph internal/repository
		B[NewMockStore] --> F[MockStore]
		F -->|Implementa| G[InsightRepository]
	end
	subgraph internal/service
		C[NewInsightService] --> H[InsightService]
		H -->|Usa| G
		H -->|Gera| I[GenerateReport]
	end
	subgraph internal/domain
		G -.->|Usa| J[Modelos: Sale, Suggestion, Complaint, UserCapture, InsightReport]
	end
	subgraph internal/transport/http
		D[NewHandler] --> E[Routes]
		E -->|Define| K[Handlers: /health, /insights/report]
		K -->|Chama| I
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

# P2 — Product Catalog API (Clean Architecture)

API REST completa para catálogo de produtos com autenticação JWT, cache Redis e Clean Architecture.

## Objetivo de aprendizado

- Clean Architecture em Go (domain, usecase, repository, handler)
- API REST com `chi` ou `gin`
- Autenticação JWT com middleware
- PostgreSQL com `sqlx`
- Cache com Redis
- Testes de integração com banco real
- Documentação automática com Swagger

## Stack

| Dependência | Uso |
|-------------|-----|
| `github.com/go-chi/chi/v5` | Router HTTP |
| `github.com/jmoiron/sqlx` | Wrapper tipado para database/sql |
| `github.com/lib/pq` | Driver PostgreSQL |
| `github.com/redis/go-redis/v9` | Cliente Redis |
| `github.com/golang-jwt/jwt/v5` | JWT tokens |
| `github.com/swaggo/swag` | Geração de Swagger |
| `github.com/stretchr/testify` | Testes |

## Arquitetura

```
internal/
├── domain/
│   ├── product.go          ← struct Product, interface ProductRepository
│   └── category.go
├── usecase/
│   ├── create_product.go   ← regra de negócio pura
│   ├── list_products.go
│   └── update_product.go
├── repository/
│   ├── postgres/
│   │   └── product_repo.go ← implementação PostgreSQL
│   └── cache/
│       └── redis_repo.go   ← cache read-through
├── handler/
│   ├── product_handler.go  ← HTTP handlers (recebe DTO, chama usecase)
│   └── auth_handler.go
├── middleware/
│   ├── auth.go             ← valida JWT, injeta userID no context
│   └── ratelimit.go
└── dto/
    ├── product_dto.go      ← request/response structs
    └── auth_dto.go

cmd/
└── api/
    └── main.go             ← monta DI, inicia servidor

docker-compose.yml           ← PostgreSQL + Redis
```

## Fluxo de dependências (Clean Architecture)

```
Handler → UseCase → Repository (interface)
                         ↑
               PostgreSQL / Redis (implementação)
```

O domínio não importa nada externo. O usecase depende apenas de interfaces.

## Como construir do zero

### 1. Setup inicial

```bash
go mod init github.com/seu-usuario/product-catalog-api
go get github.com/go-chi/chi/v5
go get github.com/jmoiron/sqlx github.com/lib/pq
go get github.com/redis/go-redis/v9
go get github.com/golang-jwt/jwt/v5
```

### 2. Definir o domínio (internal/domain/product.go)

```go
type Product struct {
    ID          string
    Name        string
    Description string
    Price       float64
    CategoryID  string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ProductRepository interface {
    Create(ctx context.Context, p Product) error
    FindByID(ctx context.Context, id string) (*Product, error)
    FindAll(ctx context.Context, filter ProductFilter) ([]Product, int, error)
    Update(ctx context.Context, p Product) error
    Delete(ctx context.Context, id string) error
}
```

### 3. Criar o usecase (internal/usecase/create_product.go)

```go
type CreateProductUseCase struct {
    repo domain.ProductRepository
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
    // validações de negócio aqui
    p := domain.Product{ID: uuid.New().String(), ...}
    return &p, uc.repo.Create(ctx, p)
}
```

### 4. Implementar o repository PostgreSQL

- Migration SQL com tabelas `products` e `categories`
- Usar `sqlx.NamedExecContext` para inserts com structs
- Paginação cursor-based: `WHERE id > $1 ORDER BY id LIMIT $2`

### 5. Implementar cache Redis

- Wrapper que implementa `ProductRepository`
- `FindByID`: tenta Redis primeiro, fallback para PostgreSQL + set no cache
- TTL de 5 minutos por produto

### 6. Criar handlers HTTP (internal/handler/product_handler.go)

```go
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var input dto.CreateProductRequest
    json.NewDecoder(r.Body).Decode(&input)
    // validar input
    result, err := h.usecase.Execute(r.Context(), input.ToUseCaseInput())
    // retornar JSON
}
```

### 7. Middleware de autenticação JWT

```go
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization") // "Bearer <token>"
            // validar JWT, extrair claims
            // ctx = context.WithValue(ctx, "userID", claims.Subject)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 8. Montar a aplicação (cmd/api/main.go)

```go
db := sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_URL")})

productRepo := postgres.NewProductRepository(db)
cachedRepo := cache.NewRedisCacheRepository(rdb, productRepo)
createUC := usecase.NewCreateProductUseCase(cachedRepo)
productHandler := handler.NewProductHandler(createUC, ...)

r := chi.NewRouter()
r.Use(middleware.Logger, middleware.Recoverer)
r.Route("/api/v1", func(r chi.Router) {
    r.Use(AuthMiddleware(os.Getenv("JWT_SECRET")))
    r.Post("/products", productHandler.Create)
    r.Get("/products", productHandler.List)
})
http.ListenAndServe(":8080", r)
```

## docker-compose.yml

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: catalog
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
    ports: ["5432:5432"]

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

### Comandos

```bash
# Subir dependências
make docker-up

# Rodar migrations
make migrate

# Iniciar API
make run

# Testar
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Notebook","price":4999.99,"category_id":"<uuid>"}'
```

### Makefile

```makefile
docker-up:
    docker compose up -d

docker-down:
    docker compose down

migrate:
    go run ./cmd/migrate/

run:
    go run ./cmd/api/

test:
    go test ./... -v -cover

test-integration:
    DATABASE_URL=postgres://admin:secret@localhost:5432/catalog go test ./... -tags=integration
```

## Desafios extras

- [ ] Rate limiting por IP com sliding window no Redis
- [ ] Upload de imagem do produto para S3 (aws-sdk-go-v2)
- [ ] Paginação cursor-based em vez de offset
- [ ] Webhook para notificar sistemas externos quando produto é criado
- [ ] GitHub Actions com testes rodando contra banco real

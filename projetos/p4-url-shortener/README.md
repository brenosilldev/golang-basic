# P4 — URL Shortener com Métricas e Alta Performance

Encurtador de URLs com contador de cliques em tempo real, expiração automática, analytics e benchmarking de performance.

## Objetivo de aprendizado

- Redis como store primário (TTL, atomic INCR, hash maps)
- Goroutines para tracking assíncrono (não bloquear o redirect)
- Atomic counters com `sync/atomic`
- Prometheus para métricas de negócio
- Benchmarking com `k6` ou `hey`
- Otimização: latência < 5ms no p99 para o redirect

## Stack

| Dependência | Uso |
|-------------|-----|
| `github.com/go-chi/chi/v5` | Router HTTP |
| `github.com/redis/go-redis/v9` | Store principal + cache |
| `github.com/lib/pq` + `sqlx` | PostgreSQL para analytics histórico |
| `github.com/prometheus/client_golang` | Métricas |

## Arquitetura

```
POST /shorten                GET /:code
     │                           │
     ▼                           ▼
┌──────────┐              ┌─────────────┐
│ Shorten  │              │  Redirect   │
│ Handler  │              │  Handler    │
└────┬─────┘              └──────┬──────┘
     │                           │
     ▼                           ▼
┌──────────────────────────────────────┐
│              Redis                   │
│  url:{code}  →  {originalURL, TTL}   │
│  clicks:{code}  →  contador INCR     │
└──────────────────────────────────────┘
                    │
                    │ async (goroutine)
                    ▼
┌──────────────────────────────────────┐
│  Analytics Writer (channel buffer)   │
│  Batches clicks → PostgreSQL         │
└──────────────────────────────────────┘
```

## Como construir do zero

### 1. Geração de código curto

```go
// Gera código base62 de 7 caracteres (62^7 = 3.5 bilhões de URLs)
func generateCode() string {
    const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, 7)
    rand.Read(b)
    for i := range b {
        b[i] = alphabet[int(b[i])%len(alphabet)]
    }
    return string(b)
}
```

### 2. Store Redis

```go
type URLStore struct {
    client *redis.Client
}

func (s *URLStore) Save(ctx context.Context, code, originalURL string, ttl time.Duration) error {
    return s.client.Set(ctx, "url:"+code, originalURL, ttl).Err()
}

func (s *URLStore) Get(ctx context.Context, code string) (string, error) {
    return s.client.Get(ctx, "url:"+code).Result()
}

func (s *URLStore) IncrClicks(ctx context.Context, code string) error {
    return s.client.Incr(ctx, "clicks:"+code).Err()
}
```

### 3. Redirect handler (caminho crítico para performance)

```go
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    code := chi.URLParam(r, "code")

    originalURL, err := h.store.Get(r.Context(), code)
    if err == redis.Nil {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        http.Error(w, "internal error", 500)
        return
    }

    // Tracking assíncrono — não bloqueia o redirect
    go h.trackClick(code, r.Header.Get("User-Agent"), r.RemoteAddr)

    http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (h *Handler) trackClick(code, ua, ip string) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    h.store.IncrClicks(ctx, code)
    h.clickCh <- ClickEvent{Code: code, UA: ua, IP: ip, At: time.Now()}
}
```

### 4. Analytics writer em background

```go
// Agrega cliques em batch para evitar write amplification no PostgreSQL
func (h *Handler) analyticsWorker(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    batch := make([]ClickEvent, 0, 100)

    for {
        select {
        case event := <-h.clickCh:
            batch = append(batch, event)
            if len(batch) >= 100 {
                h.flushAnalytics(ctx, batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                h.flushAnalytics(ctx, batch)
                batch = batch[:0]
            }
        case <-ctx.Done():
            h.flushAnalytics(ctx, batch)
            return
        }
    }
}
```

### 5. Métricas Prometheus

```go
var (
    redirectTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "urlshortener_redirects_total",
        Help: "Total de redirects por code",
    }, []string{"code"})

    redirectDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "urlshortener_redirect_duration_seconds",
        Buckets: prometheus.DefBuckets,
    }, []string{"status"})
)
```

Expor em `GET /metrics` com `promhttp.Handler()`.

### 6. Benchmarking

```bash
# Instalar hey
go install github.com/rakyll/hey@latest

# Benchmark de redirect (meta: < 5ms p99)
hey -n 100000 -c 100 http://localhost:8080/abc1234

# Output esperado:
# Latency distribution:
#   50% in 0.0012 secs
#   95% in 0.0035 secs
#   99% in 0.0048 secs  ← deve ser < 5ms
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

### docker-compose.yml

```yaml
services:
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: urlshortener
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
    ports: ["5432:5432"]
  prometheus:
    image: prom/prometheus
    volumes: ["./prometheus.yml:/etc/prometheus/prometheus.yml"]
    ports: ["9090:9090"]
  grafana:
    image: grafana/grafana
    ports: ["3000:3000"]
```

### Comandos

```bash
make docker-up
make run

# Criar URL curta
curl -X POST http://localhost:8080/shorten \
  -d '{"url":"https://github.com","ttl_hours":24}'
# {"code":"abc1234","short_url":"http://localhost:8080/abc1234"}

# Redirect
curl -L http://localhost:8080/abc1234

# Analytics
curl http://localhost:8080/analytics/abc1234

# Métricas
open http://localhost:9090  # Prometheus
open http://localhost:3000  # Grafana
```

### Makefile

```makefile
docker-up:
    docker compose up -d

run:
    go run ./cmd/

test:
    go test ./... -v -race -cover

bench:
    hey -n 50000 -c 50 http://localhost:8080/$(CODE)
```

## Desafios extras

- [ ] Vanity URLs: `/github` → github.com (usuário escolhe o código)
- [ ] QR Code gerado automaticamente para cada URL curta
- [ ] Dashboard de analytics com top 10 URLs mais acessadas
- [ ] Proteção contra abuso: rate limiting por IP com Redis
- [ ] Detectar e bloquear URLs maliciosas com Google Safe Browsing API

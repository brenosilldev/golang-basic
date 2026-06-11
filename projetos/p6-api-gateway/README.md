# P6 — API Gateway com Circuit Breaker e Service Discovery

Reverse proxy em Go que faz roteamento, retry, circuit breaking, rate limiting e service discovery. Entende o que Nginx/Envoy fazem por dentro.

## Objetivo de aprendizado

- `net/http/httputil.ReverseProxy` — como proxy funciona internamente
- Circuit breaker pattern (estados: closed, open, half-open)
- Service discovery com etcd ou arquivo de configuração dinâmica
- Rate limiting com sliding window usando Redis
- Load balancing (round-robin, least connections)
- Health checks automáticos com failover

## Stack

| Dependência | Uso |
|-------------|-----|
| `net/http/httputil` (stdlib) | Reverse proxy base |
| `github.com/redis/go-redis/v9` | Rate limiting sliding window |
| `go.etcd.io/etcd/client/v3` | Service discovery (opcional) |
| `golang.org/x/time/rate` | Token bucket rate limiter |

## Arquitetura

```
Cliente
  │
  │ HTTP Request
  ▼
┌───────────────────────────────────────────┐
│               API Gateway                 │
│                                           │
│  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │  Rate    │→ │  Router  │→ │Circuit  │ │
│  │ Limiter  │  │(por path)│  │Breaker  │ │
│  └──────────┘  └──────────┘  └────┬────┘ │
│                                   │      │
│  ┌──────────────────────────────────────┐ │
│  │        Load Balancer                │ │
│  │   Round-robin / least connections   │ │
│  └──────┬──────────┬──────────┬───────┘ │
└─────────┼──────────┼──────────┼─────────┘
          │          │          │
          ▼          ▼          ▼
      Service A  Service B  Service C
      :8081      :8082      :8083
```

## Configuração (config.yaml)

```yaml
gateway:
  port: 8080
  rate_limit:
    requests_per_second: 100
    burst: 20

routes:
  - path: /api/products
    service: product-service
    strip_prefix: true
  - path: /api/orders
    service: order-service

services:
  product-service:
    instances:
      - url: http://localhost:8081
      - url: http://localhost:8082
    health_check: /health
    circuit_breaker:
      threshold: 5      # falhas antes de abrir
      timeout: 30s      # tempo em estado open

  order-service:
    instances:
      - url: http://localhost:8083
    health_check: /health
```

## Como construir do zero

### 1. Proxy reverso base

```go
type Gateway struct {
    routes  map[string]*ServiceGroup
    config  Config
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    route := g.matchRoute(r.URL.Path)
    if route == nil {
        http.NotFound(w, r)
        return
    }

    instance, err := route.NextInstance() // load balancer
    if err != nil {
        http.Error(w, "service unavailable", 503)
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(instance.URL)
    proxy.ServeHTTP(w, r)
}
```

### 2. Circuit Breaker (state machine)

```go
type CircuitState int

const (
    StateClosed   CircuitState = iota // normal: deixa passar
    StateOpen                         // bloqueado: rejeita imediatamente
    StateHalfOpen                     // testando: deixa 1 request passar
)

type CircuitBreaker struct {
    mu          sync.Mutex
    state       CircuitState
    failures    int
    threshold   int
    lastFailure time.Time
    timeout     time.Duration
}

func (cb *CircuitBreaker) Allow() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            return true
        }
        return false
    case StateHalfOpen:
        return true
    }
    return false
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures = 0
    cb.state = StateClosed
}

func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures++
    cb.lastFailure = time.Now()
    if cb.failures >= cb.threshold {
        cb.state = StateOpen
    }
}
```

### 3. Load Balancer round-robin

```go
type ServiceGroup struct {
    instances []*Instance
    current   uint64 // atomic counter
    cb        *CircuitBreaker
}

func (sg *ServiceGroup) NextInstance() (*Instance, error) {
    if !sg.cb.Allow() {
        return nil, errors.New("circuit breaker open")
    }

    // round-robin com atomic (thread-safe sem mutex)
    idx := atomic.AddUint64(&sg.current, 1)
    instance := sg.instances[idx%uint64(len(sg.instances))]

    if !instance.Healthy {
        return sg.nextHealthyInstance()
    }
    return instance, nil
}
```

### 4. Health checks em background

```go
func (g *Gateway) startHealthChecks(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    for {
        select {
        case <-ticker.C:
            for _, group := range g.routes {
                for _, instance := range group.instances {
                    go g.checkHealth(ctx, instance)
                }
            }
        case <-ctx.Done():
            return
        }
    }
}

func (g *Gateway) checkHealth(ctx context.Context, inst *Instance) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    resp, err := http.Get(inst.URL.String() + inst.HealthPath)
    inst.Healthy = err == nil && resp.StatusCode == 200
}
```

### 5. Rate limiting com sliding window

```go
func RateLimitMiddleware(rdb *redis.Client, rps int) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := "rl:" + r.RemoteAddr
            now := time.Now().UnixMilli()
            windowStart := now - 1000 // 1 segundo

            pipe := rdb.Pipeline()
            pipe.ZRemRangeByScore(r.Context(), key, "0", strconv.FormatInt(windowStart, 10))
            pipe.ZAdd(r.Context(), key, redis.Z{Score: float64(now), Member: now})
            pipe.ZCard(r.Context(), key)
            pipe.Expire(r.Context(), key, 2*time.Second)
            results, _ := pipe.Exec(r.Context())

            count := results[2].(*redis.IntCmd).Val()
            if count > int64(rps) {
                w.WriteHeader(http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker (para Redis)

### Comandos

```bash
# Subir Redis
docker run -d -p 6379:6379 redis:alpine

# Iniciar gateway
go run ./cmd/ --config config.yaml

# Testar roteamento
curl http://localhost:8080/api/products

# Ver status dos serviços
curl http://localhost:8080/_gateway/status
```

### Makefile

```makefile
run:
    go run ./cmd/ --config config.yaml

test:
    go test ./... -v -race -cover

docker-up:
    docker run -d -p 6379:6379 redis:alpine
```

## Desafios extras

- [ ] Recarregar configuração sem reiniciar (hot reload com SIGHUP)
- [ ] Dashboard HTML com status dos serviços em tempo real (SSE)
- [ ] Retry automático com jitter em falhas 5xx
- [ ] mTLS entre gateway e serviços downstream
- [ ] Plugin system para middleware customizado (injeção via config)

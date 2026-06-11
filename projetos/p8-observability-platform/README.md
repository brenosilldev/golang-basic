# P8 — Observability Platform (Projeto Sênior)

Plataforma de observabilidade que coleta métricas e logs de múltiplas fontes, agrega, armazena e alerta em tempo real via WebSocket.

## Objetivo de aprendizado

- OpenTelemetry Go SDK (OTEL) para coleta padronizada
- gRPC para ingestão de alta performance
- Time-series storage (implementar storage simples ou integrar com InfluxDB)
- WebSocket para dashboard live sem polling
- Alertas stateful com debounce e janelas de tempo
- Processamento de 100k eventos/s com latência < 10ms no p99
- Goroutines avançadas: fan-in, fan-out, pipelines

## Stack

| Dependência | Uso |
|-------------|-----|
| `go.opentelemetry.io/otel` | SDK OTEL para instrumentação |
| `google.golang.org/grpc` | Ingestão gRPC de alta performance |
| `github.com/gorilla/websocket` | Dashboard live |
| `github.com/InfluxCommunity/influxdb3-go` | Time-series storage |
| `github.com/prometheus/client_golang` | Expose métricas próprias da plataforma |

## Arquitetura

```
Agentes / SDKs
    │ gRPC (métricas, traces, logs)
    ▼
┌──────────────────────────────────────────┐
│          Collector Service               │
│                                          │
│  gRPC Server → Ingestion Pipeline        │
│      │                                   │
│      ├── Sampler (reduz volume)          │
│      ├── Enricher (adiciona metadata)    │
│      └── Router → Storage Backends      │
└────────────────────┬─────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
  InfluxDB        AlertEngine   EventBus
  (métricas)      (regras)      (pub/sub)
                     │               │
                     ▼               ▼
               Notification      WebSocket
               (Slack/webhook)   Dashboard
```

## Pipeline de processamento (fan-in/fan-out)

```go
// Fan-in: múltiplos coletores → um canal
func MergeStreams(streams ...<-chan Metric) <-chan Metric {
    out := make(chan Metric, 1000)
    var wg sync.WaitGroup
    for _, stream := range streams {
        wg.Add(1)
        go func(s <-chan Metric) {
            defer wg.Done()
            for m := range s {
                out <- m
            }
        }(stream)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

// Fan-out: um canal → storage + alertas + websocket
func (p *Pipeline) Process(ctx context.Context, in <-chan Metric) {
    storageCh := make(chan Metric, 500)
    alertCh   := make(chan Metric, 500)
    wsCh      := make(chan Metric, 100)

    go p.storage.Write(ctx, storageCh)
    go p.alertEngine.Evaluate(ctx, alertCh)
    go p.wsHub.Broadcast(ctx, wsCh)

    for m := range in {
        storageCh <- m
        alertCh   <- m
        select {
        case wsCh <- m:
        default: // ws é best-effort, não bloqueia o pipeline
        }
    }
}
```

## Como construir do zero

### 1. Definir o schema de métricas (proto/metrics.proto)

```protobuf
syntax = "proto3";

service MetricCollector {
    rpc Send(stream MetricBatch) returns (SendResponse);
}

message MetricBatch {
    string source      = 1;
    int64  timestamp   = 2;
    repeated Metric metrics = 3;
}

message Metric {
    string name   = 1;
    double value  = 2;
    map<string, string> labels = 3;
    MetricType type = 4;
}

enum MetricType {
    GAUGE     = 0;
    COUNTER   = 1;
    HISTOGRAM = 2;
}
```

### 2. gRPC Server de ingestão

```go
type CollectorServer struct {
    pb.UnimplementedMetricCollectorServer
    pipeline *Pipeline
}

func (s *CollectorServer) Send(stream pb.MetricCollector_SendServer) error {
    for {
        batch, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.SendResponse{Received: true})
        }
        if err != nil {
            return err
        }
        for _, m := range batch.Metrics {
            s.pipeline.Ingest(Metric{
                Name:      m.Name,
                Value:     m.Value,
                Labels:    m.Labels,
                Timestamp: time.Unix(0, batch.Timestamp),
                Source:    batch.Source,
            })
        }
    }
}
```

### 3. Alert Engine com janelas de tempo

```go
type AlertRule struct {
    Name      string
    Metric    string
    Condition string // "avg > 90 for 5m"
    Threshold float64
    Window    time.Duration
}

type AlertEngine struct {
    rules   []AlertRule
    windows map[string]*SlidingWindow // metric -> janela de valores
    mu      sync.RWMutex
}

func (e *AlertEngine) Evaluate(ctx context.Context, in <-chan Metric) {
    for m := range in {
        e.mu.Lock()
        window := e.getOrCreateWindow(m.Name, 5*time.Minute)
        window.Add(m.Value, m.Timestamp)

        for _, rule := range e.rules {
            if rule.Metric == m.Name {
                avg := window.Average()
                if avg > rule.Threshold {
                    e.fire(Alert{Rule: rule.Name, Value: avg})
                }
            }
        }
        e.mu.Unlock()
    }
}
```

### 4. WebSocket Hub para dashboard live

```go
type Hub struct {
    clients   map[*Client]bool
    broadcast chan []byte
    register  chan *Client
    mu        sync.RWMutex
}

func (h *Hub) Run(ctx context.Context) {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
        case data := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- data:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        case <-ctx.Done():
            return
        }
    }
}
```

### 5. Benchmarking de 100k eventos/s

```go
func BenchmarkPipeline_Throughput(b *testing.B) {
    pipeline := NewPipeline(...)
    metric := Metric{Name: "cpu.usage", Value: 42.5}

    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pipeline.Ingest(metric)
    }
    // meta: > 100k ops/s (b.N/tempo < 10μs por op)
}
```

### 6. Agente OTEL para instrumentar aplicações

```go
// Exemplo de agente que coleta métricas do sistema e envia via gRPC
func main() {
    conn, _ := grpc.Dial("collector:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
    client := pb.NewMetricCollectorClient(conn)
    stream, _ := client.Send(context.Background())

    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        var mem runtime.MemStats
        runtime.ReadMemStats(&mem)

        stream.Send(&pb.MetricBatch{
            Source:    "my-app",
            Timestamp: time.Now().UnixNano(),
            Metrics: []*pb.Metric{
                {Name: "heap_alloc_bytes", Value: float64(mem.HeapAlloc)},
                {Name: "goroutine_count", Value: float64(runtime.NumGoroutine())},
            },
        })
    }
}
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`

### Comandos

```bash
# Gerar código gRPC
make proto

# Subir infraestrutura
make docker-up

# Iniciar collector
make run-collector

# Iniciar agente em uma aplicação qualquer
make run-agent

# Abrir dashboard
open http://localhost:3000
```

### Makefile

```makefile
proto:
    protoc --go_out=. --go-grpc_out=. proto/metrics.proto

docker-up:
    docker compose up -d

run-collector:
    go run ./cmd/collector/

run-agent:
    go run ./cmd/agent/

test:
    go test ./... -v -race -cover

bench:
    go test ./... -bench=. -benchmem -benchtime=10s
```

## Desafios extras

- [ ] Retenção automática: dados > 30 dias são downsampled ou deletados
- [ ] Correlação de traces com métricas (trace_id nos labels)
- [ ] Multi-tenant: cada organização vê apenas seus dados
- [ ] Exportar para Prometheus / Grafana como backend alternativo
- [ ] Anomaly detection simples com Z-score em tempo real

# P3 — Worker Pool com Processamento de Jobs

Sistema que recebe jobs via API HTTP, processa com um pool limitado de goroutines e reporta o status em tempo real.

## Objetivo de aprendizado

- Goroutines e channels (buffered e unbuffered)
- WaitGroup, Mutex e atomic operations
- Bounded concurrency (sem spawnar goroutine ilimitada)
- Graceful shutdown com context cancelation
- Backpressure e retry com exponential backoff
- Circuit breaker pattern

## Stack

| Dependência | Uso |
|-------------|-----|
| `net/http` (stdlib) | Servidor HTTP |
| `database/sql` + `modernc.org/sqlite` | Persistência de estado dos jobs |
| `golang.org/x/sync/errgroup` | Gerenciar grupo de goroutines com erro |

## Arquitetura

```
                   HTTP POST /jobs
                        │
                        ▼
              ┌─────────────────┐
              │   Job Queue     │  channel buffered (capacidade N)
              │  chan Job       │
              └────────┬────────┘
                       │ distribui para
          ┌────────────┼────────────┐
          ▼            ▼            ▼
      Worker 1     Worker 2     Worker 3   (goroutines fixas)
          │            │            │
          └────────────┴────────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  Job Store      │  SQLite / in-memory
              │  (status, errs) │
              └─────────────────┘
                       │
                       ▼
            GET /jobs/:id/status
```

## Como construir do zero

### 1. Definir o Job

```go
type JobStatus string

const (
    StatusQueued     JobStatus = "queued"
    StatusProcessing JobStatus = "processing"
    StatusDone       JobStatus = "done"
    StatusFailed     JobStatus = "failed"
)

type Job struct {
    ID        string
    Payload   json.RawMessage
    Status    JobStatus
    Retries   int
    Error     string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 2. Criar o Dispatcher

O dispatcher é o coração do sistema: aceita jobs e os distribui para workers.

```go
type Dispatcher struct {
    queue      chan Job          // jobs aguardando processamento
    results    chan JobResult    // resultados dos workers
    numWorkers int
    store      JobStore
    processor  JobProcessor     // interface com Execute(ctx, job) error
    wg         sync.WaitGroup
}

func (d *Dispatcher) Start(ctx context.Context) {
    // Inicia numWorkers goroutines
    for i := 0; i < d.numWorkers; i++ {
        d.wg.Add(1)
        go d.worker(ctx, i)
    }
    // Goroutine separada para processar resultados
    go d.processResults(ctx)
}

func (d *Dispatcher) worker(ctx context.Context, id int) {
    defer d.wg.Done()
    for {
        select {
        case job, ok := <-d.queue:
            if !ok {
                return // canal fechado, worker encerra
            }
            d.processJob(ctx, job)
        case <-ctx.Done():
            return
        }
    }
}
```

### 3. Lógica de retry com exponential backoff

```go
func (d *Dispatcher) processJob(ctx context.Context, job Job) {
    const maxRetries = 3
    var err error

    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            backoff := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
            select {
            case <-time.After(backoff):
            case <-ctx.Done():
                return
            }
        }
        err = d.processor.Execute(ctx, job)
        if err == nil {
            d.results <- JobResult{JobID: job.ID, Status: StatusDone}
            return
        }
    }
    d.results <- JobResult{JobID: job.ID, Status: StatusFailed, Error: err.Error()}
}
```

### 4. Graceful Shutdown

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Captura SIGTERM/SIGINT
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

    go func() {
        <-sigCh
        cancel() // sinaliza para todos os workers pararem
    }()

    dispatcher.Start(ctx)

    <-ctx.Done()
    close(dispatcher.queue) // fecha o canal após cancelamento
    dispatcher.wg.Wait()    // espera workers terminarem jobs em andamento
    log.Println("Shutdown completo, nenhum job perdido")
}
```

### 5. API HTTP

```
POST /jobs          → enfileira um job, retorna {"id": "...", "status": "queued"}
GET  /jobs/:id      → retorna status atual do job
GET  /jobs          → lista todos com filtro por status
GET  /metrics       → workers ativos, jobs pendentes, taxa de sucesso
```

### 6. Testar concorrência

```go
func TestWorkerPool_ProcessesConcurrently(t *testing.T) {
    // Criar dispatcher com 3 workers
    // Enviar 100 jobs
    // Verificar que todos foram processados
    // Verificar que no máximo 3 estavam em processamento simultaneamente
    // Medir tempo: deve ser ~jobs/workers * tempo_por_job
}
```

## Como rodar

### Pré-requisitos

- Go 1.22+

### Comandos

```bash
# Iniciar servidor (8 workers, fila de 100)
NUM_WORKERS=8 QUEUE_SIZE=100 go run ./cmd/

# Enviar job
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"type":"email","payload":{"to":"user@example.com","subject":"Hello"}}'

# Verificar status
curl http://localhost:8080/jobs/<id>

# Ver métricas
curl http://localhost:8080/metrics
```

### Makefile

```makefile
run:
    go run ./cmd/

test:
    go test ./... -v -race -cover

bench:
    go test ./... -bench=. -benchmem
```

> **Importante:** sempre rode com `-race` para detectar race conditions.

## Desafios extras

- [ ] Prioridade de jobs (fila de alta prioridade tem preferência)
- [ ] Dead letter queue para jobs que falharam após max retries
- [ ] Persistência no PostgreSQL para sobreviver a reinicializações
- [ ] Dashboard em tempo real via Server-Sent Events (SSE)
- [ ] Circuit breaker: se > 50% dos jobs falham, parar de processar por X segundos

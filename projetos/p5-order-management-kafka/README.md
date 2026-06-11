# P5 — Order Management System com Kafka (Microsserviços)

Sistema de pedidos event-driven com 3 microsserviços independentes comunicando via Kafka. Projeto principal de portfólio.

## Objetivo de aprendizado

- Arquitetura event-driven com Kafka
- Saga pattern para transações distribuídas
- Idempotência e at-least-once delivery
- gRPC entre serviços internos
- Dead letter queue para falhas
- Docker Compose orquestrando múltiplos serviços
- Consistência eventual na prática

## Stack

| Dependência | Uso |
|-------------|-----|
| `github.com/segmentio/kafka-go` | Produtor/consumidor Kafka |
| `google.golang.org/grpc` | Comunicação síncrona entre serviços |
| `github.com/jmoiron/sqlx` + driver pq | PostgreSQL por serviço |
| `github.com/google/uuid` | IDs únicos para idempotência |
| `github.com/go-chi/chi/v5` | HTTP gateway externo |

## Arquitetura

```
Cliente
  │
  │ HTTP POST /orders
  ▼
┌──────────────────┐
│  order-service   │ → publica OrderCreated no Kafka
│  PostgreSQL: orders│
└──────────────────┘
         │ Kafka topic: orders.created
         ▼
┌──────────────────┐
│ payment-service  │ → processa pagamento
│ PostgreSQL: payments│→ publica PaymentProcessed ou PaymentFailed
└──────────────────┘
    │                │
    │ Kafka:         │ Kafka:
    │ payment.done   │ payment.failed
    ▼                ▼
┌──────────────────┐  ┌──────────────────┐
│notification-svc  │  │  order-service   │ → atualiza status para FAILED
│ envia email/hook │  │  saga compensação│
└──────────────────┘  └──────────────────┘
```

## Eventos (Kafka Topics)

| Topic | Produtor | Consumidores |
|-------|----------|--------------|
| `orders.created` | order-service | payment-service |
| `payment.processed` | payment-service | order-service, notification-service |
| `payment.failed` | payment-service | order-service (compensação) |
| `orders.completed` | order-service | notification-service |
| `orders.dlq` | qualquer serviço | monitoramento |

## Como construir do zero

### 1. Estrutura de pastas

```
order-management/
├── docker-compose.yml
├── services/
│   ├── order-service/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── domain/order.go
│   │   │   ├── handler/order_handler.go
│   │   │   ├── producer/kafka_producer.go
│   │   │   └── consumer/payment_consumer.go
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── payment-service/
│   │   └── ... (mesma estrutura)
│   └── notification-service/
│       └── ...
└── proto/
    └── order.proto    ← definição gRPC compartilhada
```

### 2. Definir eventos (contratos Kafka)

```go
// shared/events/events.go (copiado para cada serviço ou em módulo compartilhado)

type OrderCreatedEvent struct {
    EventID    string    `json:"event_id"`    // UUID para idempotência
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    Items      []Item    `json:"items"`
    TotalCents int64     `json:"total_cents"`
    CreatedAt  time.Time `json:"created_at"`
}

type PaymentProcessedEvent struct {
    EventID     string    `json:"event_id"`
    OrderID     string    `json:"order_id"`
    PaymentID   string    `json:"payment_id"`
    Status      string    `json:"status"` // "approved" | "rejected"
    ProcessedAt time.Time `json:"processed_at"`
}
```

### 3. Produtor Kafka (order-service)

```go
type OrderProducer struct {
    writer *kafka.Writer
}

func (p *OrderProducer) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    payload, _ := json.Marshal(event)
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(event.OrderID), // garante ordering por order
        Value: payload,
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte("OrderCreated")},
            {Key: "event_id", Value: []byte(event.EventID)},
        },
    })
}
```

### 4. Consumidor com idempotência (payment-service)

```go
func (c *OrderConsumer) Start(ctx context.Context) {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            break
        }

        var event OrderCreatedEvent
        json.Unmarshal(msg.Value, &event)

        // Idempotência: checar se event_id já foi processado
        if c.store.EventProcessed(ctx, event.EventID) {
            continue // já processado, ignorar (at-least-once delivery)
        }

        err = c.processPayment(ctx, event)
        if err != nil {
            c.publishToDeadLetterQueue(ctx, msg, err)
            continue
        }

        c.store.MarkEventProcessed(ctx, event.EventID)
        c.reader.CommitMessages(ctx, msg)
    }
}
```

### 5. Saga de compensação

Se o pagamento falhar, order-service precisa reverter o pedido:

```go
// consumer que ouve payment.failed no order-service
func (c *PaymentFailedConsumer) Handle(ctx context.Context, event PaymentFailedEvent) error {
    return c.orderRepo.UpdateStatus(ctx, event.OrderID, StatusCancelled)
}
```

### 6. Dead Letter Queue

Mensagens que falharam após N tentativas vão para `orders.dlq`:

```go
func (c *Consumer) publishToDeadLetterQueue(ctx context.Context, original kafka.Message, err error) {
    dlqWriter.WriteMessages(ctx, kafka.Message{
        Value: original.Value,
        Headers: append(original.Headers, kafka.Header{
            Key: "error", Value: []byte(err.Error()),
        }),
    })
}
```

## docker-compose.yml (resumido)

```yaml
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
    ports: ["9092:9092"]

  order-service:
    build: ./services/order-service
    environment:
      KAFKA_BROKERS: kafka:9092
      DATABASE_URL: postgres://...
    ports: ["8080:8080"]
    depends_on: [kafka, postgres-orders]

  payment-service:
    build: ./services/payment-service
    environment:
      KAFKA_BROKERS: kafka:9092
    depends_on: [kafka]

  notification-service:
    build: ./services/notification-service
    depends_on: [kafka]

  postgres-orders:
    image: postgres:16-alpine

  postgres-payments:
    image: postgres:16-alpine
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

### Comandos

```bash
# Subir toda a infraestrutura
make docker-up

# Criar pedido
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "123",
    "items": [{"product_id": "abc", "quantity": 2, "price_cents": 9999}]
  }'
# {"order_id": "...", "status": "pending"}

# Acompanhar status
curl http://localhost:8080/orders/<order_id>
# {"status": "completed"} ou {"status": "cancelled"}

# Ver logs dos serviços
docker compose logs -f payment-service
```

### Makefile

```makefile
docker-up:
    docker compose up -d --build

docker-down:
    docker compose down -v

test:
    go test ./services/... -v -cover

test-integration:
    docker compose -f docker-compose.test.yml up --abort-on-container-exit
```

## Desafios extras

- [ ] Kafka UI (conduktor/kafka-ui) para visualizar topics e mensagens
- [ ] Outbox pattern para garantir atomicidade entre DB e Kafka
- [ ] Circuit breaker no payment-service para falhas externas
- [ ] Tracing distribuído com OpenTelemetry conectando todos os serviços
- [ ] Testes de contrato entre serviços com Pact

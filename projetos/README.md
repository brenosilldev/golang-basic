# Projetos GoExpert → Engenheiro Sênior

Roadmap de projetos práticos baseados no curso GoExpert FullCycle, ordenados por complexidade crescente.

## Projetos

| # | Projeto | Nível | Conceitos principais |
|---|---------|-------|----------------------|
| P1 | [Task Manager CLI](./p1-task-manager-cli/) | Iniciante | CLI, SQLite, structs, interfaces |
| P2 | [Product Catalog API](./p2-product-catalog-api/) | Intermediário | REST, Clean Arch, JWT, PostgreSQL |
| P3 | [Worker Pool](./p3-worker-pool/) | Intermediário+ | Goroutines, channels, concorrência |
| P4 | [URL Shortener](./p4-url-shortener/) | Intermediário+ | Redis, Prometheus, performance |
| P5 | [Order Management + Kafka](./p5-order-management-kafka/) | Avançado | Microsserviços, Kafka, saga pattern |
| P6 | [API Gateway](./p6-api-gateway/) | Avançado | Proxy reverso, circuit breaker, service discovery |
| P7 | [Deploy CLI](./p7-deploy-cli/) | Avançado | Docker SDK, YAML, DevOps |
| P8 | [Observability Platform](./p8-observability-platform/) | Sênior | OpenTelemetry, gRPC, alta performance |

## Estratégia

```
Semanas 1-4   → P1 + "Let's Go Further" (Alex Edwards)
Semanas 5-10  → P2 + módulo de concorrência do GoExpert
Semanas 11-16 → P3 + MIT 6.824 Distributed Systems (YouTube)
Semanas 17-24 → P5 — projeto principal de portfólio
Paralelo      → P7 em sprints de fim de semana
```

## Padrão de qualidade para cada projeto

- `README.md` com arquitetura e como rodar
- `Makefile` com `make run`, `make test`, `make docker-up`
- Cobertura de testes > 70%
- GitHub Actions para CI

## Recursos de estudo

| Recurso | Tipo | Por que |
|---------|------|---------|
| Let's Go Further — Alex Edwards | Livro | Melhor referência para APIs avançadas em Go |
| MIT 6.824 Distributed Systems | Curso gratuito | O mais desafiador — Raft, MapReduce, Spanner |
| Designing Data-Intensive Applications | Livro | Base teórica de sistemas distribuídos |
| FullCycle 3.0 | Curso | Docker, K8s, DDD, arquitetura hexagonal |
| System Design Interview Vol. 1 e 2 | Livros | Desenho de sistemas para entrevistas sênior |

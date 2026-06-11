# P7 — Deploy CLI (mini-Helm em Go)

Ferramenta CLI que lê um manifesto YAML e gerencia o ciclo de vida de containers Docker: build, deploy, health check e rollback automático.

## Objetivo de aprendizado

- Docker SDK para Go (`docker/docker`)
- Parsing de YAML e template engine
- CLI avançada com cobra (subcomandos aninhados)
- Lifecycle management de containers
- Rollback automático baseado em health check
- Goroutines para operações paralelas de deploy

## Stack

| Dependência | Uso |
|-------------|-----|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/docker/docker/client` | Docker Engine API |
| `gopkg.in/yaml.v3` | Parsing de manifesto |
| `github.com/docker/docker/api/types` | Tipos Docker |
| `golang.org/x/term` | Progress bars no terminal |

## Manifesto (deploy.yaml)

```yaml
app: my-api
version: "1.2.0"

services:
  api:
    image: "my-api:{{ .Version }}"
    build:
      context: .
      dockerfile: Dockerfile
    replicas: 3
    ports:
      - "8080:8080"
    env:
      DATABASE_URL: "${DATABASE_URL}"
      LOG_LEVEL: "info"
    health_check:
      path: /health
      interval: 5s
      timeout: 2s
      retries: 3
    resources:
      memory: "256m"
      cpu: "0.5"

  worker:
    image: "my-worker:{{ .Version }}"
    replicas: 2
    env:
      QUEUE_URL: "${QUEUE_URL}"
```

## Arquitetura

```
CLI Command
    │
    ▼
┌────────────────────────────────────┐
│         Deploy Engine              │
│                                    │
│  Parse YAML → Template Render      │
│       │                            │
│       ▼                            │
│  Plan (o que vai mudar)            │
│       │                            │
│       ▼                            │
│  Execute Deploy (Docker API)       │
│  ┌─────────────────────────────┐   │
│  │  pull image                 │   │
│  │  stop old containers        │   │
│  │  start new containers       │   │
│  │  wait health checks         │   │
│  │  remove old containers      │   │
│  └─────────────────────────────┘   │
│       │              │             │
│  health OK?      health FAIL?      │
│       │              │             │
│   Done ✓         Rollback ↩        │
└────────────────────────────────────┘
```

## Como construir do zero

### 1. Definir o manifesto (internal/manifest/manifest.go)

```go
type Manifest struct {
    App      string              `yaml:"app"`
    Version  string              `yaml:"version"`
    Services map[string]*Service `yaml:"services"`
}

type Service struct {
    Image      string            `yaml:"image"`
    Build      *BuildConfig      `yaml:"build"`
    Replicas   int               `yaml:"replicas"`
    Ports      []string          `yaml:"ports"`
    Env        map[string]string `yaml:"env"`
    HealthCheck *HealthCheck     `yaml:"health_check"`
}

func Load(path string, vars map[string]string) (*Manifest, error) {
    raw, _ := os.ReadFile(path)

    // Template rendering: substitui {{ .Version }} etc.
    tmpl, _ := template.New("manifest").Parse(string(raw))
    var buf bytes.Buffer
    tmpl.Execute(&buf, vars)

    var m Manifest
    yaml.Unmarshal(buf.Bytes(), &m)
    return &m, nil
}
```

### 2. Cliente Docker (internal/docker/client.go)

```go
type DockerClient struct {
    cli *client.Client
}

func (d *DockerClient) PullImage(ctx context.Context, image string) error {
    reader, err := d.cli.ImagePull(ctx, image, image.PullOptions{})
    if err != nil {
        return err
    }
    defer reader.Close()
    io.Copy(os.Stdout, reader) // mostra progresso no terminal
    return nil
}

func (d *DockerClient) StartContainer(ctx context.Context, cfg ContainerConfig) (string, error) {
    resp, err := d.cli.ContainerCreate(ctx,
        &container.Config{
            Image: cfg.Image,
            Env:   cfg.EnvSlice(),
        },
        &container.HostConfig{
            PortBindings: cfg.PortBindings(),
            Resources: container.Resources{
                Memory:   cfg.MemoryBytes(),
                NanoCPUs: cfg.CPUNanoCores(),
            },
        },
        nil, nil, cfg.Name,
    )
    if err != nil {
        return "", err
    }
    return resp.ID, d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
}
```

### 3. Health check (internal/deploy/healthcheck.go)

```go
func WaitHealthy(ctx context.Context, containerPort string, hc HealthCheck) error {
    url := fmt.Sprintf("http://localhost:%s%s", containerPort, hc.Path)
    deadline := time.Now().Add(time.Duration(hc.Retries) * hc.Interval)

    for time.Now().Before(deadline) {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(hc.Interval):
            resp, err := http.Get(url)
            if err == nil && resp.StatusCode == 200 {
                return nil
            }
        }
    }
    return fmt.Errorf("health check falhou após %d tentativas", hc.Retries)
}
```

### 4. Deploy com rollback (internal/deploy/deploy.go)

```go
func (d *Deployer) Deploy(ctx context.Context, svc *Service) error {
    // 1. Pull nova imagem
    if err := d.docker.PullImage(ctx, svc.Image); err != nil {
        return err
    }

    // 2. Guardar IDs dos containers atuais para rollback
    oldContainers, _ := d.docker.ListContainers(ctx, svc.Name)

    // 3. Subir novos containers
    var newContainers []string
    for i := 0; i < svc.Replicas; i++ {
        id, err := d.docker.StartContainer(ctx, ContainerConfig{
            Name:  fmt.Sprintf("%s-%s-%d", svc.Name, svc.Version, i),
            Image: svc.Image,
        })
        if err != nil {
            d.rollback(ctx, newContainers, oldContainers)
            return err
        }
        newContainers = append(newContainers, id)
    }

    // 4. Health check
    if err := WaitHealthy(ctx, svc.Ports[0], *svc.HealthCheck); err != nil {
        fmt.Println("Health check falhou, executando rollback...")
        d.rollback(ctx, newContainers, oldContainers)
        return err
    }

    // 5. Remover containers antigos
    for _, id := range oldContainers {
        d.docker.StopContainer(ctx, id)
        d.docker.RemoveContainer(ctx, id)
    }
    return nil
}
```

### 5. Subcomandos CLI

```bash
# Subcomandos disponíveis
deploy up --file deploy.yaml          # sobe serviços
deploy down --file deploy.yaml        # para e remove containers
deploy status --file deploy.yaml      # lista containers e health
deploy rollback --file deploy.yaml    # volta para versão anterior
deploy plan --file deploy.yaml        # mostra o que vai mudar (dry-run)
deploy logs --service api --tail 100  # logs de um serviço
```

## Como rodar

### Pré-requisitos

- Go 1.22+
- Docker rodando localmente

### Comandos

```bash
# Compilar
go build -o deploy ./cmd/

# Ver plano antes de aplicar
./deploy plan --file deploy.yaml

# Fazer deploy
./deploy up --file deploy.yaml --version 1.2.0

# Verificar status
./deploy status --file deploy.yaml

# Rollback se algo der errado
./deploy rollback --file deploy.yaml
```

### Makefile

```makefile
build:
    go build -o deploy ./cmd/

install:
    go install ./cmd/

test:
    go test ./... -v -cover

lint:
    golangci-lint run
```

## Desafios extras

- [ ] Suporte a Docker Compose files além do formato próprio
- [ ] Deploy em Kubernetes usando `k8s.io/client-go`
- [ ] Notificações de deploy via Slack webhook
- [ ] Histórico de deploys com possibilidade de rollback para versão específica
- [ ] Secrets management integrado com Vault ou AWS SSM

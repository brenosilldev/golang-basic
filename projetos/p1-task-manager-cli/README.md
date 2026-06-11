# P1 — Task Manager CLI

Ferramenta de linha de comando para gerenciar tarefas com persistência local em SQLite.

## Objetivo de aprendizado

- Construir CLI com subcomandos usando `cobra`
- Persistência de dados com SQLite em Go
- Structs e interfaces para modelagem do domínio
- Testes unitários com `testing` e `testify`
- Serialização/deserialização (JSON e CSV)

## Stack

| Dependência | Uso |
|-------------|-----|
| `github.com/spf13/cobra` | Framework de CLI com subcomandos |
| `modernc.org/sqlite` | SQLite puro Go (sem cgo) |
| `github.com/stretchr/testify` | Assertions em testes |

## Arquitetura

```
cmd/
├── main.go              ← entrypoint
├── root.go              ← comando raiz e configurações globais
├── add.go               ← task add "titulo" --tag work
├── list.go              ← task list --status pending --tag work
├── done.go              ← task done <id>
├── delete.go            ← task delete <id>
└── export.go            ← task export --format csv

internal/
├── task/
│   ├── task.go          ← struct Task, interface Repository
│   └── task_test.go
└── storage/
    ├── sqlite.go        ← implementação SQLite do Repository
    └── sqlite_test.go
```

## Como construir do zero

### 1. Inicializar o projeto

```bash
mkdir task-manager-cli && cd task-manager-cli
go mod init github.com/seu-usuario/task-manager-cli
go get github.com/spf13/cobra@latest
go get modernc.org/sqlite@latest
go get github.com/stretchr/testify@latest
```

### 2. Definir o domínio (internal/task/task.go)

```go
type Status string

const (
    StatusPending  Status = "pending"
    StatusDone     Status = "done"
)

type Task struct {
    ID        int
    Title     string
    Tags      []string
    Status    Status
    CreatedAt time.Time
}

type Repository interface {
    Save(t Task) (Task, error)
    FindAll(filter Filter) ([]Task, error)
    FindByID(id int) (Task, error)
    UpdateStatus(id int, status Status) error
    Delete(id int) error
}
```

### 3. Implementar o storage SQLite (internal/storage/sqlite.go)

- Criar tabela `tasks` com colunas: id, title, tags (JSON), status, created_at
- Implementar cada método da interface `Repository`
- Abrir conexão com `sql.Open("sqlite", "~/.tasks.db")`

### 4. Criar os subcomandos cobra (cmd/)

Cada arquivo implementa um `cobra.Command` e injeta o `Repository`:

```go
// cmd/add.go
var addCmd = &cobra.Command{
    Use:   "add [titulo]",
    Short: "Adiciona uma nova tarefa",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        tag, _ := cmd.Flags().GetString("tag")
        // chamar repository.Save(...)
    },
}
```

### 5. Conectar tudo em cmd/root.go

```go
func Execute(repo task.Repository) {
    rootCmd.AddCommand(addCmd, listCmd, doneCmd, deleteCmd, exportCmd)
    rootCmd.Execute()
}
```

### 6. Implementar exportação CSV/JSON (cmd/export.go)

- `--format csv`: usar `encoding/csv`
- `--format json`: usar `encoding/json`

### 7. Escrever testes

- `internal/task/task_test.go`: testar lógica de negócio (filtros, validações)
- `internal/storage/sqlite_test.go`: testar CRUD com banco em memória (`:memory:`)

## Como rodar

### Pré-requisitos

- Go 1.22+

### Comandos

```bash
# Compilar
go build -o task ./cmd/

# Adicionar tarefa
./task add "Estudar goroutines" --tag study

# Listar tarefas pendentes
./task list --status pending

# Marcar como concluída
./task done 1

# Exportar para CSV
./task export --format csv > tarefas.csv
```

### Makefile

```makefile
build:
    go build -o task ./cmd/

test:
    go test ./... -v -cover

run:
    go run ./cmd/ list
```

## Desafios extras

- [ ] Adicionar prioridade (low/medium/high) e ordenar por ela
- [ ] Suporte a due date com alertas de vencimento
- [ ] Subcomando `task stats` mostrando quantidade por status/tag
- [ ] Colorir output no terminal com `github.com/fatih/color`

# Conceitos fundamentais — Go Expert

---

## Concorrência

### Goroutines e o scheduler M:P:G

**O que é:** Goroutines são threads leves gerenciadas pelo runtime Go, não pelo SO. O scheduler usa três entidades: M (thread de SO), P (processador lógico, controla a fila de goroutines) e G (goroutine). GOMAXPROCS define quantos P existem simultâneamente.

**Como funciona em Go:** Cada P tem uma run queue local. O scheduler é preemptivo desde Go 1.14 (via sinais SIGURG), então loops sem syscalls não bloqueiam mais os outros. Goroutines migram entre P quando um P fica ocioso (work stealing). Criar uma goroutine custa ~2KB de stack inicial — pode existir milhões.

**Por que importa na prática:** Entender M:P:G explica por que `GOMAXPROCS=1` ainda paraleliza I/O, por que goroutines não são threads, e por que `runtime.LockOSThread()` é necessário em integrações com libs C que exigem thread-affinity.

**Armadilha comum:** Assumir que goroutines rodam em paralelo quando `GOMAXPROCS=1` — elas são concorrentes, não paralelas. O scheduler coopera em pontos de preempção (chamadas de função, syscalls, channel ops).

**Exercício mental:** Você tem 10.000 goroutines fazendo CPU-bound work em um servidor com 4 cores. Qual `GOMAXPROCS` faz sentido e por quê? E se o trabalho for I/O-bound?

---

### Channels (unbuffered vs buffered)

**O que é:** Channels são o mecanismo de comunicação entre goroutines. Unbuffered: o send bloqueia até alguém receber (sincronização explícita). Buffered: o send só bloqueia quando o buffer está cheio (desacoplamento parcial).

**Como funciona em Go:** `make(chan T)` cria unbuffered. `make(chan T, n)` cria buffered com capacidade n. Internamente, um channel é uma estrutura com mutex, circular buffer, e filas de goroutines suspensas esperando para enviar ou receber.

```go
ch := make(chan int, 3) // buffer de 3
ch <- 1  // não bloqueia — há espaço
ch <- 2  // não bloqueia
ch <- 3  // não bloqueia
ch <- 4  // BLOQUEIA — buffer cheio
```

**Por que importa na prática:** A escolha entre buffered e unbuffered muda a semântica de sincronização. Unbuffered garante que o sender sabe que o receiver já pegou a mensagem. Buffered desacopla produtor e consumidor, mas esconde backpressure se o tamanho for arbitrário.

**Armadilha comum:** Usar buffer grande demais para "evitar bloqueios" — isso elimina backpressure e esconde o fato de que o consumidor está lento. O sistema cresce sem limites e explode na memória.

**Exercício mental:** Em um pipeline A → B → C, onde B é lento, qual estratégia de buffering mantém A e C trabalhando sem acumular dados indefinidamente?

---

### Select statement

**O que é:** `select` permite aguardar múltiplas operações de channel simultaneamente. Quando mais de um case está pronto, Go escolhe um aleatoriamente (distribuição uniforme, não FIFO).

**Como funciona em Go:** O runtime registra todos os cases, tenta cada um sem bloquear, e se nenhum estiver pronto suspende a goroutine até um case ficar disponível. `default` torna o select não-bloqueante.

```go
select {
case msg := <-ch1:
    process(msg)
case ch2 <- result:
    // enviou
case <-ctx.Done():
    return ctx.Err()
default:
    // nenhum canal pronto agora
}
```

**Por que importa na prática:** Padrão fundamental para timeout, cancelamento via context, multiplexação de fontes, e heartbeats. Sem `select`, cada goroutine bloquearia em um único canal.

**Armadilha comum:** Colocar `default` em loops de polling — cria busy-wait que consome 100% de CPU sem fazer trabalho real. Use `time.After` ou `time.Ticker` junto com o select.

**Exercício mental:** Você quer processar mensagens de dois channels mas dar prioridade ao channel de controle sobre o de dados. O select aleatório é suficiente? Como garantiria prioridade real?

---

### sync.Mutex e sync.RWMutex

**O que é:** `sync.Mutex` é um mutex binário — lock/unlock. `sync.RWMutex` distingue leitores e escritores: múltiplos leitores simultâneos são permitidos, mas um escritor tem exclusividade total.

**Como funciona em Go:** `Mutex` usa operações atômicas internamente com fallback para syscall quando há contenção. `RWMutex` mantém contador de leitores ativos; um escritor bloqueia novos leitores e espera os existentes terminarem.

**Por que importa na prática:** RWMutex é a escolha certa para caches e estruturas com leitura dominante (ex: 99% reads, 1% writes). Usar Mutex nesse caso serializa operações que poderiam ser paralelas.

**Armadilha comum:** Passar `Mutex` por valor — isso copia o estado interno e cria dois mutexes independentes. Sempre passe ponteiro, ou use o mutex como campo de uma struct e passe a struct por ponteiro.

**Exercício mental:** Você tem um map em memória lido por 100 goroutines e escrito por 1. Qual a diferença de throughput entre `sync.Mutex` e `sync.RWMutex` sob carga? E comparado com `sync.Map`?

---

### sync.WaitGroup

**O que é:** Contador que permite uma goroutine esperar um grupo de outras terminar. `Add(n)` incrementa, `Done()` decrementa (equivale a `Add(-1)`), `Wait()` bloqueia até o contador chegar a zero.

**Como funciona em Go:** Implementado com contador atômico e semáforo. `Wait()` dorme até o contador zerar. Quando a última goroutine chama `Done()`, todos os waiters são acordados.

**Por que importa na prática:** O padrão básico de fan-out em Go. Sem WaitGroup, você precisaria de um channel extra de sinalização e coordenação manual de contagem.

**Armadilha comum:** Chamar `Add` dentro da goroutine lançada em vez de antes de lançá-la. A goroutine pode terminar antes do `Add` executar, fazendo o `Wait` retornar cedo.

**Exercício mental:** Você lança 5 goroutines com `go func()` e chama `wg.Add(1)` dentro de cada uma. O que acontece se o scheduler demorar para iniciar as goroutines?

---

### sync/atomic

**O que é:** Operações atômicas em nível de instrução de CPU — leitura, escrita, CAS (Compare-And-Swap), Add. Não usam lock, portanto são mais rápidas que Mutex para operações simples em valores únicos.

**Como funciona em Go:** `atomic.AddInt64`, `atomic.LoadInt64`, `atomic.StoreInt64`, `atomic.CompareAndSwapInt64`. No Go 1.19+, `atomic.Int64` é o wrapper tipado recomendado. Internamente mapeiam para instruções LOCK do x86 ou equivalentes.

**Por que importa na prática:** Contadores de métricas (requests, erros, bytes), flags booleanos de estado (running/stopped), e ponteiros que trocam atomicamente (hot reload de config). Qualquer Mutex em torno de um único inteiro provavelmente deve ser atomic.

**Armadilha comum:** Usar atomic para múltiplas variáveis que precisam ser consistentes entre si — atomic garante atomicidade por variável, não entre variáveis. Para isso, Mutex ou uma struct atômica com CAS são necessários.

**Exercício mental:** Você tem `hits int64` e `misses int64` em um cache. Pode usar atomic.Add em ambos independentemente? E se precisar calcular hit rate lendo os dois ao mesmo tempo?

---

### Context (cancelamento, timeout, propagação)

**O que é:** `context.Context` é a forma idiomática Go de propagar cancelamento, deadlines e valores por uma árvore de chamadas. Quando um contexto pai é cancelado, todos os filhos são cancelados automaticamente.

**Como funciona em Go:** `context.WithCancel` retorna um filho e uma função `cancel()`. `context.WithTimeout`/`WithDeadline` adicionam deadline. O channel `ctx.Done()` fecha quando cancelado. `ctx.Err()` retorna a razão.

**Por que importa na prática:** Toda operação que faz I/O (banco, HTTP, gRPC, Redis) deve receber `ctx` como primeiro argumento. Sem context, você não consegue cancelar uma query quando o cliente desconecta — a operação continua consumindo recursos.

**Armadilha comum:** Guardar context em structs como campo. Context é de escopo de request/operation, não de vida útil do objeto. Passá-lo como argumento de função é a forma correta.

**Exercício mental:** Uma requisição HTTP recebe contexto com timeout de 5s. Ela faz 3 chamadas sequenciais a serviços externos, cada uma podendo levar até 3s. O que acontece com as chamadas 2 e 3 se a 1 demorar 4s? Sem context propagado, o que mudaria?

---

### Race conditions e como detectar com go test -race

**O que é:** Uma race condition ocorre quando duas goroutines acessam a mesma memória concorrentemente e pelo menos uma delas escreve, sem sincronização. O resultado depende do timing — não-determinístico.

**Como funciona em Go:** O flag `-race` ativa o race detector do compilador, baseado no ThreadSanitizer. Ele instrumenta cada acesso a memória e rastreia a ordem de acesso via vector clocks. É ~5-10x mais lento mas detecta races em tempo de execução.

```
go test -race ./...
go run -race main.go
```

**Por que importa na prática:** Races causam corrupção de dados silenciosa, crashes intermitentes e comportamentos impossíveis de reproduzir. O detector elimina toda uma categoria de bugs que seriam horas de debugging em produção.

**Armadilha comum:** Só rodar `-race` em testes unitários simples. Races em código de produção frequentemente só se manifestam sob carga. Rode testes de integração e benchmarks com `-race` também.

**Exercício mental:** Você tem um map global `cache map[string]int` lido e escrito por goroutines. Os testes unitários passam sem `-race`. Quando você adiciona `-race`, ele falha. Por que os testes sem `-race` não detectaram o problema?

---

### Fan-out e fan-in patterns

**O que é:** Fan-out distribui trabalho de um producer para múltiplos workers concorrentes. Fan-in agrega resultados de múltiplos producers em um único channel de saída.

**Como funciona em Go:** Fan-out: um channel de entrada, N goroutines lendo dele. Fan-in: N goroutines escrevendo em channels separados, uma goroutine de merge usando `select` para combinar em um único canal de saída.

**Por que importa na prática:** O par fan-out/fan-in é a base de qualquer pipeline de processamento paralelo — extração de dados, transformações independentes, chamadas a serviços externos em paralelo.

**Armadilha comum:** No fan-in, fazer N goroutines escrevendo num único channel sem fechá-lo corretamente. O receptor precisa saber quando todos os producers terminaram — use WaitGroup + goroutine de fechamento do canal agregador.

**Exercício mental:** Você faz 5 chamadas de API externas em paralelo (fan-out) e agrega os resultados (fan-in). Uma das APIs demora 30s. Como garante que um timeout de 5s cancela todas as chamadas, incluindo as que ainda não começaram?

---

### Backpressure via channels

**O que é:** Backpressure é o mecanismo pelo qual um consumidor lento sinaliza ao produtor para desacelerar, evitando acúmulo ilimitado de dados em memória.

**Como funciona em Go:** Em Go, backpressure é natural em channels bufferizados: quando o buffer enche, o send bloqueia o produtor. Isso propaga a lentidão para cima na cadeia. Channels unbuffered têm backpressure máxima — o produtor só avança quando o consumidor aceita.

**Por que importa na prática:** Sem backpressure, um produtor rápido e consumidor lento acumula dados indefinidamente, causando OOM ou latência crescente. Backpressure é o que transforma um sistema que parece funcionar em carga baixa em um que sobrevive em produção.

**Armadilha comum:** Usar goroutine que sempre lê do canal e joga num slice em memória para "não bloquear o produtor" — isso simplesmente move o buffer ilimitado para o heap e o problema persiste.

**Exercício mental:** Seu producer gera 10.000 eventos/s e o consumer processa 1.000/s. Com buffer de 100, em quanto tempo o sistema satura? O que acontece com a latência do producer quando satura?

---

### Worker pool pattern

**O que é:** Um número fixo de goroutines (workers) processam jobs de uma fila compartilhada. O pool tem tamanho estático ou dinâmico. Impede que carga pico crie goroutines ilimitadas.

**Como funciona em Go:** Canal de jobs como fila, N goroutines lendo desse canal. O caller envia jobs ao canal. Quando o canal enche, o caller bloqueia (backpressure). Fechar o canal de jobs encerra os workers após drenagem.

**Por que importa na prática:** Essencial para qualquer operação custosa (queries de banco, chamadas HTTP, processamento de arquivos) onde não se quer N goroutines para N requisições simultâneas — o overhead de contexto e os recursos (conexões de banco, descritores) explodem.

**Armadilha comum:** Criar o pool mas nunca fechar o canal de jobs ou esperar os workers terminarem — goroutines ficam suspensas para sempre e o programa termina antes de processar todos os jobs.

**Exercício mental:** Você tem um pool de 10 workers e recebe 1000 jobs de 100ms cada. Qual a taxa máxima de throughput? E se um worker travar, os outros 9 continuam? Como você detectaria o worker travado?

---

### Circuit breaker

**O que é:** Padrão de resiliência que interrompe chamadas a um serviço quando a taxa de falhas supera um limiar, evitando cascata de erros. Três estados: Closed (normal), Open (falha rápida), Half-Open (testando recuperação).

**Como funciona em Go:** O circuit breaker mantém janela de contagem de sucessos e falhas. Ao atingir o threshold de falha, transiciona para Open e retorna erro imediatamente sem fazer a chamada real. Após timeout, vai para Half-Open e permite uma chamada de prova.

**Por que importa na prática:** Sem circuit breaker, quando o serviço B cai, o serviço A acumula threads/goroutines esperando timeout, o que degrada A também. O circuit breaker rompe a cascata em milissegundos.

**Armadilha comum:** Configurar o threshold muito sensível (ex: 1 falha abre o circuito) — transient errors normais abrem o circuit breaker e causam indisponibilidade desnecessária. Use janela deslizante e percentual mínimo de chamadas antes de avaliar.

**Exercício mental:** Seu serviço tem timeout de 30s por chamada ao downstream. Com 100 requisições simultâneas, quanto tempo leva para esgotar o pool de conexões sem circuit breaker? E com circuit breaker de 5s de timeout aberto?

---

## Memória e Runtime

### Stack vs heap em Go

**O que é:** Stack é memória privada de cada goroutine, gerenciada automaticamente pelo runtime (push/pop de frames). Heap é memória compartilhada gerenciada pelo GC. Alocações no stack são praticamente gratuitas; heap tem overhead de GC.

**Como funciona em Go:** Go stacks começam pequenos (~2KB) e crescem sob demanda (até o limite configurado). O compilador decide onde alocar cada variável via escape analysis. Variáveis que "escapam" do frame da função vão para o heap.

**Por que importa na prática:** Código de hot path que aloca muito no heap aumenta pressão no GC e pode causar latência de GC visível (stop-the-world, mesmo que curto). Entender stack vs heap é o primeiro passo para otimização de alocações.

**Armadilha comum:** Retornar ponteiro para variável local e assumir que está no stack — se ela escapa, está no heap. Use `go build -gcflags="-m"` para ver as decisões de escape analysis do compilador.

**Exercício mental:** Uma função cria um slice de 1MB e retorna um ponteiro para ele. Onde esse slice mora? E se a função retornasse o slice por valor (não ponteiro)?

---

### Escape analysis

**O que é:** Análise estática do compilador que determina se uma variável pode ser alocada no stack ou precisa ir para o heap. Uma variável "escapa" quando sua vida útil excede o frame da função onde foi criada.

**Como funciona em Go:** O compilador rastreia se um ponteiro para a variável é devolvido ao caller, armazenado numa interface, enviado para channel, ou capturado por closure que pode sobreviver à função. Se qualquer um desses ocorre, a variável escapa para o heap.

```
go build -gcflags="-m=2" ./...
# ./main.go:10:6: moved to heap: x
```

**Por que importa na prática:** Interfaces são um dos maiores gatilhos de escape. Um `fmt.Sprintf` aloca no heap porque os argumentos são passados como `interface{}`. Em hot paths (ex: serialização), isso pode ser significativo.

**Armadilha comum:** Passar structs pequenas por ponteiro "por performance" — se o ponteiro causa escape para o heap, o overhead de GC pode superar o custo de copiar a struct. Structs pequenas frequentemente ficam mais rápidas passadas por valor.

**Exercício mental:** Você tem uma função que recebe `interface{}` e faz type assertion. O valor passado escapa para o heap? E se você usasse uma interface tipada como `io.Reader` em vez de `interface{}`?

---

### Garbage collector do Go (tri-color mark-and-sweep)

**O que é:** O GC do Go é concorrente, não-compactador, tri-color mark-and-sweep. "Tri-color" significa que objetos são brancos (não visitados), cinzas (visitados, filhos pendentes) ou pretos (visitados, filhos também visitados). O GC roda majoritariamente concorrente com a aplicação.

**Como funciona em Go:** O GC tem duas fases de stop-the-world curtas (ativar write barriers e finalizar marcação) e uma fase de marcação concorrente longa. A write barrier garante que objetos criados durante a marcação não sejam perdidos (invariante tri-color).

**Por que importa na prática:** O GC do Go é otimizado para latência baixa, não throughput de GC. A meta padrão é `GOGC=100` (heap dobra antes de coletar). Aumentar `GOGC` reduz frequência de GC mas aumenta uso de memória; diminuir faz o oposto.

**Armadilha comum:** Alocar muitos objetos pequenos de vida curta em hot paths — isso aumenta o trabalho do marcador mesmo que a memória total seja pequena. Pooling com `sync.Pool` é a solução padrão.

**Exercício mental:** Seu serviço tem picos de latência a cada ~2 minutos, mesmo sem aumento de carga. O que você suspeitaria primeiro? Como confirmaria?

---

### Snapshots e por que reduzem GC pressure

**O que é:** Em vez de alocar objetos individuais frequentemente, você tira um "snapshot" do estado atual do sistema — cópia imutável de uma estrutura de dados — e usa essa cópia até a próxima atualização. O GC vê poucos objetos de longa vida em vez de muitos de curta.

**Como funciona em Go:** Padrão comum com `sync.Map` ou `atomic.Pointer[T]`: writer constrói novo estado completo, swapa atomicamente o ponteiro. Readers leem o ponteiro atual sem lock. Nenhum objeto antigo é modificado — apenas substituído.

**Por que importa na prática:** Para caches de configuração, routing tables, feature flags — qualquer estado lido com frequência mas escrito raramente. O padrão troca custo de GC por custo de cópia completa a cada write.

**Armadilha comum:** Usar snapshot para dados que mudam com alta frequência — a cópia completa a cada write pode ser mais cara que o custo de GC evitado. Meça antes de aplicar.

**Exercício mental:** Você tem uma lista de 10.000 regras de roteamento que muda 1 vez por minuto mas é lida 1 milhão de vezes por minuto. Snapshot ou RWMutex? Por quê?

---

### pprof: CPU profile, heap profile, goroutine dump

**O que é:** `pprof` é a ferramenta de profiling do Go. CPU profile mostra onde o programa gasta tempo. Heap profile mostra onde aloca memória (allocs vs in-use). Goroutine dump lista todas as goroutines ativas com suas stacks.

**Como funciona em Go:** Importe `_ "net/http/pprof"` e registre o handler. Acesse `/debug/pprof/profile?seconds=30` para CPU, `/debug/pprof/heap` para heap, `/debug/pprof/goroutine` para goroutines. Use `go tool pprof` para análise interativa e flame graphs.

**Por que importa na prática:** Sem profiling, otimizações são chute. pprof revela com precisão onde o tempo é gasto e onde a memória é alocada. Flame graphs tornam óbvio o que não é óbvio olhando o código.

**Armadilha comum:** Expor o endpoint `/debug/pprof` sem autenticação em produção — o profiling tem overhead e expõe informações internas. Use firewall ou autenticação básica.

**Exercício mental:** Você ativa CPU profiling por 30s e vê que 40% do tempo está em `runtime.mallocgc`. O que isso indica? Qual o próximo passo de diagnóstico?

---

### Benchmarks com testing.B e como medir P50/P95/P99

**O que é:** Benchmarks em Go usam `testing.B`, são executados com `go test -bench=.` e medem operações por segundo (ns/op). Por padrão mostram média. Para percentis (P50/P95/P99), use o pacote `benchstat` ou `histogram`.

**Como funciona em Go:** `b.N` é ajustado automaticamente pelo framework para que o benchmark rode por tempo suficiente. Use `b.ResetTimer()` após setup. Use `b.RunParallel` para benchmarks concorrentes. Para percentis: capture tempos individuais num slice e calcule com `sort`.

**Por que importa na prática:** Média esconde outliers. P99 é o que o usuário mais lento experimenta. Em sistemas distribuídos, P99 de um serviço se torna P50 quando há 100 dependências (tail amplification). Medir apenas média é enganoso.

**Armadilha comum:** Benchmarking sem `-benchmem` — você vê ns/op mas não sabe quantas alocações estão acontecendo. Uma função rápida que aloca muito pode ser mais lenta em produção por causa do GC.

**Exercício mental:** Seu benchmark mostra 500 ns/op e 2 allocs/op. Outro candidato mostra 800 ns/op e 0 allocs/op. Qual você escolhe para um hot path chamado 100.000 vezes por segundo? Por quê?

---

## Banco de Dados e Consistência

### ACID vs BASE

**O que é:** ACID (Atomicity, Consistency, Isolation, Durability) é o modelo de bancos relacionais: transações são tudo-ou-nada, isoladas umas das outras, e duráveis. BASE (Basically Available, Soft state, Eventually consistent) descreve bancos distribuídos que priorizam disponibilidade sobre consistência forte.

**Como funciona em Go:** Bancos ACID são acessados via `database/sql` com `db.BeginTx` para transações explícitas. Bancos BASE (Cassandra, DynamoDB) têm SDKs próprios onde consistência é configurada por operação (eventual vs strong read).

**Por que importa na prática:** A escolha define os garantias do seu sistema. ACID simplifica a lógica de negócio (você não precisa lidar com estados intermediários visíveis). BASE permite escala horizontal mas transfere a complexidade de consistência para a aplicação.

**Armadilha comum:** Misturar semântica ACID e BASE no mesmo fluxo — ex: gravar num banco relacional e num Kafka na mesma "transação" sem Outbox pattern. Se um falhar e o outro não, você tem inconsistência.

**Exercício mental:** Você tem um sistema de reservas de hotel. O inventário fica num banco ACID. Os emails de confirmação são enviados via fila. O banco confirma mas a fila cai. O usuário recebeu a reserva mas não o email. Isso é aceitável? E se fosse ao contrário?

---

### Transações distribuídas: por que 2PC não escala

**O que é:** Two-Phase Commit (2PC) é um protocolo que coordena commit ou rollback em múltiplos serviços/bancos. Fase 1 (Prepare): coordinator pergunta a todos se podem commitar. Fase 2 (Commit/Abort): coordinator ordena commit se todos disseram sim.

**Como funciona em Go:** 2PC em Go seria implementado via coordinator que faz chamadas RPC para os participantes. Na prática, raramente implementado diretamente — frameworks como XA o fazem, mas com overhead enorme.

**Por que importa na prática:** 2PC não escala porque o coordinator é um ponto único de falha e os participantes ficam bloqueados esperando a decisão do coordinator (locks distribuídos). Se o coordinator cai entre as fases, o sistema entra em limbo. Por isso, sistemas modernos usam Sagas.

**Armadilha comum:** Tentar emular 2PC com chamadas HTTP síncronas entre serviços — isso é 2PC sem as garantias de recovery, pegando o pior dos dois mundos.

**Exercício mental:** Em 2PC, o coordinator envia "Prepare" para 3 participantes. Dois respondem "Yes". O coordinator cai antes de enviar "Commit". Os dois participantes ficam em que estado? Por quanto tempo?

---

### Outbox pattern

**O que é:** Padrão que garante que um evento seja publicado num message broker se e somente se a operação no banco de dados for commitada. Salva o evento numa tabela "outbox" na mesma transação do dado, depois um poller lê e publica.

**Como funciona em Go:** Na mesma transação SQL: `INSERT INTO orders (...)` + `INSERT INTO outbox (event_type, payload, status='pending')`. Um processo separado (poller) lê registros `status='pending'` e publica no RabbitMQ/Kafka, depois marca como `status='published'`.

**Por que importa na prática:** Resolve o problema dual-write: sem Outbox, você faz INSERT no banco e depois publica no broker — se o processo cair entre os dois, a publicação se perde. O Outbox transforma o broker em "best-effort" com retry garantido.

**Armadilha comum:** O poller publicar e marcar como published em duas operações não-atômicas — se cair no meio, o evento é publicado duas vezes. Seus consumers precisam ser idempotentes para lidar com isso.

**Exercício mental:** Seu Outbox poller publica o evento no RabbitMQ e depois faz UPDATE outbox SET status='published'. O processo cai depois do publish mas antes do UPDATE. O que acontece quando reinicia?

---

### Idempotência e deduplication

**O que é:** Uma operação é idempotente se executá-la múltiplas vezes produz o mesmo resultado de executá-la uma vez. Deduplication é a técnica de detectar e descartar mensagens duplicadas antes de processá-las.

**Como funciona em Go:** Deduplication via tabela `processed_events (event_id PRIMARY KEY, processed_at)`. Antes de processar, tenta `INSERT INTO processed_events`. Se viola unique constraint, a mensagem já foi processada — descarte. O INSERT e o processamento ficam na mesma transação.

**Por que importa na prática:** At-least-once delivery garante que você receberá a mensagem, mas pode recebê-la mais de uma vez. Idempotência é o que transforma at-least-once em exactly-once do ponto de vista do negócio.

**Armadilha comum:** Fazer idempotência apenas na camada de aplicação sem persistir o evento processado — se o processo reinicia, a memória é perdida e o evento é reprocessado.

**Exercício mental:** Você processa um pagamento de R$100. O consumer cai depois de debitar a conta mas antes de marcar o evento como processado. Na reentrega, o evento chega de novo. Sem idempotência, o que acontece? Com idempotência via deduplication table, o que acontece?

---

### At-least-once vs exactly-once delivery

**O que é:** At-least-once: a mensagem é entregue ao menos uma vez, podendo haver duplicatas. Exactly-once: a mensagem é processada exatamente uma vez, sem duplicatas e sem perdas. At-most-once: pode ser perdida, nunca duplicada.

**Como funciona em Go:** At-least-once: consumer faz ack depois de processar com sucesso; se cair antes do ack, o broker reentregar. Exactly-once verdadeiro exige transações distribuídas ou deduplication. Kafka tem exactly-once dentro do cluster via transações de produtor.

**Por que importa na prática:** Exactly-once é a garantia mais cara. Na prática, at-least-once + idempotência dá semântica equivalente a exactly-once com custo muito menor. A maioria dos sistemas usa essa combinação.

**Armadilha comum:** Assumir que o broker garante exactly-once sozinho — mesmo brokers que afirmam isso (Kafka) têm limitações de escopo. A idempotência do consumer é sempre necessária para garantias ponta-a-ponta.

**Exercício mental:** Você tem um consumer que debita um valor de uma conta. Para ser idempotente, o que você precisa persistir além do próprio débito?

---

### Concorrência otimista (optimistic locking)

**O que é:** Estratégia de controle de concorrência que não bloqueia registros. Cada registro tem uma versão. Ao atualizar, você verifica se a versão ainda é a mesma que quando leu. Se outro processo atualizou no meio tempo, o UPDATE afeta 0 linhas — você detecta o conflito e pode tentar novamente.

**Como funciona em Go:** `UPDATE accounts SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3`. Se `RowsAffected() == 0`, houve conflito. Retorne erro de conflito ao caller para retry ou propagação.

**Por que importa na prática:** Locking pessimista (`SELECT FOR UPDATE`) serializa todas as operações naquele registro, criando contention em tabelas quentes. Optimistic locking permite paralelismo e só falha quando há conflito real — ideal quando conflitos são raros.

**Armadilha comum:** Fazer retry ilimitado em caso de conflito — em alta contenção, todos os retries também conflitam e você tem livelock. Adicione backoff exponencial com jitter e limite de tentativas.

**Exercício mental:** 100 goroutines tentam debitar de uma conta simultaneamente. Com pessimistic locking, qual a taxa de sucesso por segundo? Com optimistic locking e retry, qual o comportamento esperado?

---

### SQLC: geração de código a partir de queries SQL

**O que é:** SQLC é uma ferramenta que gera código Go type-safe a partir de queries SQL escritas manualmente. Você escreve o SQL, o SQLC gera as funções Go correspondentes com tipos corretos. Elimina ORMs e reflection em runtime.

**Como funciona em Go:** Você define queries em arquivos `.sql` com anotações (`-- name: GetUser :one`), executa `sqlc generate`, e recebe structs e funções tipadas. A geração usa o schema do banco para inferir tipos. Zero reflection em runtime.

**Por que importa na prática:** ORMs adicionam overhead de abstração e geram SQL imprevisível. SQLC te dá controle total sobre o SQL (importante para performance) com type-safety em tempo de compilação. Se o schema mudar e quebrar uma query, o código não compila.

**Armadilha comum:** Esquecer de regenerar o código após alterar as queries ou o schema — o código gerado fica defasado e o compilador não detecta porque o arquivo antigo ainda compila. Integre `sqlc generate` no processo de build.

**Exercício mental:** Você muda o tipo de uma coluna de `INT` para `BIGINT` no schema. Com ORM, quando você descobriria? Com SQLC, quando você descobriria?

---

## Arquitetura

### Clean Architecture em Go: camadas, regra de dependência

**O que é:** Organização do código em camadas concêntricas onde dependências só apontam para dentro. De fora para dentro: Frameworks/Drivers → Adapters → Use Cases → Entities. O domínio não conhece o banco, o HTTP, nem nenhuma infraestrutura.

**Como funciona em Go:** Entities em `internal/domain/`, Use Cases em `internal/usecase/`, Adapters em `internal/infra/` (implementações de interfaces do domínio). O use case depende de interfaces definidas no domínio; o adapter implementa essas interfaces.

**Por que importa na prática:** Facilita troca de tecnologia (swap de banco, de broker) sem alterar lógica de negócio. Facilita testes unitários do domínio sem infraestrutura real. Evita que complexidade acidental (frameworks) contamine complexidade essencial (negócio).

**Armadilha comum:** Criar interfaces no package de infraestrutura e depender delas no domínio — inverte a dependência no sentido errado. Interfaces de repositório devem morar no domínio ou no use case, não na infra.

**Exercício mental:** Seu use case `ProcessPayment` precisa salvar no banco e enviar email. Como você estrutura as interfaces? Quem define `PaymentRepository`? Quem implementa?

---

### Interfaces pequenas e composição

**O que é:** Em Go, interfaces são implicitamente satisfeitas e devem ter o menor número possível de métodos. A interface padrão tem 1-2 métodos. Interfaces maiores são compostas de interfaces menores.

**Como funciona em Go:** `io.Reader` tem 1 método. `io.ReadCloser` compõe `io.Reader` e `io.Closer`. Se sua função precisa apenas ler, aceite `io.Reader`, não `io.ReadCloser` — você aceita mais implementações e é mais fácil de testar (mock de 1 método vs 2).

**Por que importa na prática:** Interfaces grandes criam acoplamento desnecessário. Se você define `UserRepository` com 15 métodos, um use case que só chama `FindByID` ainda depende de todos os 15 para fins de mock/teste. Separe por capability.

**Armadilha comum:** Definir interfaces com muitos métodos "para ser completo" — isso é o Java-ism que Go ativamente desencoraja. Defina a interface mínima que o consumidor precisa, não a máxima que o tipo pode oferecer.

**Exercício mental:** Você tem `type Storage interface { Save(Entity) error; FindByID(id string) (Entity, error); FindAll() ([]Entity, error); Delete(id string) error }`. Um use case precisa apenas de `Save`. Como você refatora?

---

### Dependency Injection manual vs Google Wire

**O que é:** DI é o padrão de passar dependências para um objeto em vez de ele criá-las internamente. Manual: você escreve funções `New*` que recebem dependências como parâmetros. Wire: ferramenta que gera o código de wiring automaticamente via análise estática.

**Como funciona em Go:** DI manual: `func NewUserUseCase(repo UserRepository, email EmailSender) *UserUseCase`. Wire: você define Providers e Injectors em arquivos `wire.go`, executa `wire gen`, e ele gera o código de composição.

**Por que importa na prática:** Projetos pequenos/médios funcionam bem com DI manual — é explícito e sem magia. Projetos com dezenas de dependências tornam o bootstrap de `main.go` caótico; Wire automatiza o wiring e detecta dependências faltando em tempo de compilação.

**Armadilha comum:** Com Wire, adicionar um Provider mas esquecer de incluí-lo no ProviderSet — Wire não compila e a mensagem de erro pode ser confusa. Leia a mensagem com cuidado para identificar qual binding está faltando.

**Exercício mental:** Seu `main.go` tem 80 linhas de wiring manual. Quando vale a pena migrar para Wire? Quais são os critérios objetivos?

---

### Unit of Work pattern

**O que é:** Unit of Work (UoW) rastreia todas as mudanças feitas a objetos de domínio durante uma operação e persiste tudo numa única transação ao final. Evita que múltiplos repositórios façam commits independentes.

**Como funciona em Go:** Uma struct `UnitOfWork` wrapa uma `*sql.Tx`. Os repositórios recebem a transação em vez de uma conexão. O use case opera nos repositórios, e ao final chama `uow.Commit()` ou `uow.Rollback()`.

**Por que importa na prática:** Sem UoW, um use case que salva em duas tabelas pode commitar a primeira e falhar na segunda, deixando o banco inconsistente. UoW garante que ou tudo commita ou nada commita.

**Armadilha comum:** Criar UoW que abre transação no construtor — isso mantém transações abertas por tempo indeterminado se a construção e o uso forem em momentos distintos. Abra a transação o mais tarde possível e feche o mais cedo possível.

**Exercício mental:** Você tem `OrderRepository` e `InventoryRepository`. Um pedido cria um registro em orders e decrementa o inventário. Como Unit of Work garante atomicidade aqui?

---

### Repository pattern

**O que é:** Abstração que encapsula a lógica de acesso a dados e expõe uma coleção orientada a objetos. O domínio vê uma coleção em memória; a infraestrutura faz as queries reais. O domínio não sabe se os dados vêm de PostgreSQL, Redis ou memória.

**Como funciona em Go:** Interface definida no domínio: `type UserRepository interface { Save(ctx, *User) error; FindByID(ctx, string) (*User, error) }`. Implementação em `internal/infra/postgres/user_repository.go`. Testes usam implementação in-memory.

**Por que importa na prática:** Torna o domínio testável sem banco real. Permite trocar PostgreSQL por MySQL sem alterar use cases. Centraliza queries em um lugar só, evitando SQL espalhado pelo código.

**Armadilha comum:** Repositório com métodos como `FindByNameAndAgeGreaterThanAndStatus` — você está vazando critérios de query específicos pela interface. Use Specification pattern ou Query Objects para queries complexas.

**Exercício mental:** Seu repositório tem `FindActiveUsersByCity(city string)`. Amanhã você precisa de `FindActiveUsersByCityAndAge`. E depois de mais 5 variações. Como você evita uma explosão de métodos no repositório?

---

### Event Dispatcher e handlers desacoplados

**O que é:** Pub/sub dentro do processo: um emissor publica um evento de domínio sem saber quem vai processar. Handlers se registram para tipos de evento. O dispatcher faz o roteamento.

**Como funciona em Go:** Dispatcher mantém `map[string][]EventHandler`. `dispatcher.Register("OrderPlaced", sendEmailHandler)`. `dispatcher.Dispatch(ctx, event)` chama todos os handlers registrados para aquele tipo, síncronos ou assíncronos.

**Por que importa na prática:** Desacopla o domínio dos efeitos colaterais. `PlaceOrder` não precisa saber que vai enviar email, criar notificação e atualizar analytics — ele só emite `OrderPlaced`. Novos handlers são adicionados sem tocar no domínio.

**Armadilha comum:** Dispatcher síncrono onde um handler lento bloqueia o use case inteiro. Em produção, handlers de efeitos colaterais (email, notificação) devem ser assíncronos — ou o dispatcher os executa em goroutines separadas, ou o evento vai para uma fila real.

**Exercício mental:** Um use case emite `UserCreated`. Três handlers se registram: enviar email de boas-vindas, criar perfil de analytics, e notificar admin. Um deles sempre falha. Como o dispatcher deve tratar essa falha?

---

## Eventos e Mensageria

### Event Sourcing: estado derivado de eventos

**O que é:** Em vez de persistir o estado atual de um objeto, você persiste a sequência de eventos que causaram esse estado. O estado atual é reconstruído fazendo replay dos eventos. O log de eventos é a fonte da verdade — o estado atual é uma projeção derivada.

**Como funciona em Go:** O Event Store é append-only: `INSERT INTO events (stream_id, version, type, payload)`. Para reconstituir um aggregate: `SELECT * FROM events WHERE stream_id = $1 ORDER BY version`. Aplica cada evento ao aggregate em ordem.

**Por que importa na prática:** Você ganha histórico completo de auditoria gratuitamente. Pode reconstruir o estado em qualquer ponto no tempo. Pode criar projeções para diferentes casos de uso. Útil em domínios financeiros, inventário, e qualquer contexto onde "por que ficou assim?" importa.

**Armadilha comum:** Event store que permite UPDATE ou DELETE de eventos — isso destrói a imutabilidade que é a garantia fundamental do Event Sourcing. O event store deve ser genuinamente append-only no nível do banco.

**Exercício mental:** Você tem 50.000 eventos para uma conta bancária. Reconstituir o estado atual leva 5 segundos. Como você resolve isso sem abandonar Event Sourcing?

---

### CQRS: separar modelo de escrita e leitura

**O que é:** Command Query Responsibility Segregation. Modelos de escrita (commands) e leitura (queries) são separados. O modelo de escrita é otimizado para consistência e invariantes de negócio. O modelo de leitura é otimizado para queries específicas da UI.

**Como funciona em Go:** Commands passam pelo use case e alteram o aggregate. Quando o aggregate emite evento, uma projeção assíncrona atualiza uma tabela de leitura desnormalizada. Queries leem direto da tabela de leitura — sem joins complexos, sem reconstituição de aggregate.

**Por que importa na prática:** Em sistemas com padrões de leitura muito diferentes dos de escrita (ex: write por ID, read por múltiplos filtros e agregações), CQRS permite otimizar cada modelo independentemente. A leitura pode ser um banco diferente, uma view materializada, ou um Redis.

**Armadilha comum:** Aplicar CQRS a domínios simples — o overhead de sincronização entre write e read model é significativo. CQRS é justificado quando a complexidade de read ≠ complexidade de write, não como padrão default.

**Exercício mental:** Você tem um e-commerce. O modelo de escrita processa um pedido por vez. O modelo de leitura precisa mostrar "pedidos do último mês agrupados por categoria com total gasto". Como os dois modelos seriam estruturados?

---

### Aggregate e fronteiras de consistência

**O que é:** Aggregate é um cluster de objetos de domínio tratado como uma unidade de consistência. Toda mudança passa pela raiz do aggregate (Aggregate Root). Nenhum objeto externo tem referência direta aos filhos internos — apenas à raiz.

**Como funciona em Go:** `Order` (raiz) contém `[]OrderItem`. Para adicionar item: `order.AddItem(item)` — nunca `order.Items = append(order.Items, item)` externamente. A raiz valida invariantes de negócio (ex: total máximo, limite de itens).

**Por que importa na prática:** A fronteira do aggregate define a fronteira de transação. Duas operações em aggregates diferentes podem ocorrer em paralelo. Operações no mesmo aggregate são serializadas. Aggregates muito grandes criam contention; muito pequenos criam transações distribuídas.

**Armadilha comum:** Aggregate que referencia outros aggregates por valor (objeto embutido) em vez de por ID — isso viola a fronteira e cria dependências de carregamento. Referencie outros aggregates apenas por ID.

**Exercício mental:** Você tem `Order` e `Customer`. Um pedido pertence a um cliente. Deve `Order` conter `Customer` ou `CustomerID`? Quais são as implicações de cada escolha para consistência e performance?

---

### Saga pattern: orquestração vs coreografia

**O que é:** Saga é um padrão para transações distribuídas onde cada serviço executa sua parte local e publica eventos. Se um passo falha, os passos anteriores executam transações compensatórias. Orquestração: um coordinator central dirige os passos. Coreografia: serviços reagem a eventos sem coordinator central.

**Como funciona em Go:** Orquestração: `SagaOrchestrator` tem máquina de estados, chama cada serviço e trata sucesso/falha. Coreografia: cada serviço ouve eventos do anterior, processa e emite próximo evento. Falha emite evento de compensação.

**Por que importa na prática:** Sagas substituem 2PC com eventual consistency. Orquestração centraliza a lógica da saga (mais fácil de debugar, single ponto de falha). Coreografia distribui a responsabilidade (mais resiliente, mais difícil de rastrear o fluxo).

**Armadilha comum:** Compensação que falha — você precisa de retry e idempotência nas compensações também. Se a compensação falha após retry, você precisar de alerta e intervenção manual (Dead Letter Queue para sagas).

**Exercício mental:** Uma saga tem 4 passos: reservar estoque, processar pagamento, confirmar pedido, enviar email. O passo 3 falha. Quais compensações precisam ser executadas, em que ordem, e o que acontece se o pagamento já foi debitado mas a compensação do pagamento falha?

---

### RabbitMQ: exchanges, queues, bindings, ack/nack

**O que é:** No RabbitMQ, producers publicam em Exchanges, não diretamente em queues. Exchanges rotacionam mensagens para queues via bindings baseados em routing keys. Tipos de exchange: direct, fanout, topic, headers.

**Como funciona em Go:** Com `amqp091-go`: declare exchange, declare queue, faça binding. Publish com routing key. Consumer recebe `Delivery`, processa, chama `d.Ack(false)` em sucesso ou `d.Nack(false, requeue)` em falha. `autoAck: false` é mandatório para at-least-once.

**Por que importa na prática:** Exchanges desacoplam producer de queue. Um producer publica em `orders.topic`; múltiplas filas podem se bindar para diferentes routing keys (`orders.created`, `orders.cancelled`) sem o producer saber.

**Armadilha comum:** `autoAck: true` — a mensagem é removida da fila assim que entregue ao consumer, antes de ser processada. Se o consumer cai durante o processamento, a mensagem é perdida. Sempre use `autoAck: false` e faça ack manual.

**Exercício mental:** Você tem 3 consumers (email, analytics, audit) que precisam receber todos os eventos de `OrderCreated`. Você cria 1 queue ou 3? Que tipo de exchange você usa e por quê?

---

### Dead Letter Queue

**O que é:** Uma DLQ é uma fila especial que recebe mensagens que não puderam ser processadas com sucesso após N tentativas, ou que expiraram (TTL), ou que foram rejeitadas (nack sem requeue). Permite inspeção e reprocessamento manual.

**Como funciona em Go:** Configure `x-dead-letter-exchange` e `x-dead-letter-routing-key` na declaração da fila original. Quando uma mensagem é nack'd sem requeue, RabbitMQ a envia automaticamente para a DL exchange. Uma fila separada se binda à DL exchange para inspeção.

**Por que importa na prática:** Sem DLQ, mensagens que causam erro repetido bloqueiam o consumer num loop infinito ou são perdidas. Com DLQ, o fluxo principal continua e mensagens problemáticas ficam preservadas para diagnóstico e retry manual.

**Armadilha comum:** DLQ sem monitoramento — mensagens se acumulam silenciosamente. Configure alerta quando o tamanho da DLQ cresce, e um processo periódico de análise de mensagens mortas.

**Exercício mental:** Uma mensagem na DLQ contém dados de pagamento. Você corrigiu o bug que causava a falha. Como você reprocessa as mensagens da DLQ com segurança, garantindo idempotência?

---

## Comunicação entre Serviços

### gRPC: HTTP/2 multiplexing, Protobuf, streams

**O que é:** gRPC é um framework RPC que usa HTTP/2 para transporte (multiplexing de múltiplas chamadas numa única conexão TCP) e Protobuf para serialização (binário, ~10x menor que JSON). Suporta 4 tipos de chamada: unary, server-streaming, client-streaming, bidirectional-streaming.

**Como funciona em Go:** Define o serviço em `.proto`, gera código com `protoc`. O servidor implementa a interface gerada. O cliente usa o stub gerado. Uma única conexão gRPC pode ter milhares de chamadas simultâneas via HTTP/2 streams.

**Por que importa na prática:** Latência menor que REST/JSON em chamadas internas. Contrato fortemente tipado via `.proto` evita erros de campo. Streaming bidirecional é ideal para subscriptions, progress updates e comunicação em tempo real.

**Armadilha comum:** Criar uma nova conexão gRPC por request — conexões gRPC são custosas de criar (TLS handshake, HTTP/2 setup). Use uma pool ou compartilhe uma única `*grpc.ClientConn` thread-safe entre todos os requests.

**Exercício mental:** Você tem 10.000 RPCs simultâneos para o mesmo servidor. Com REST/HTTP1.1, quantas conexões TCP? Com gRPC/HTTP2, quantas?

---

### Resolver e Balancer API do gRPC-Go

**O que é:** A Resolver API permite que o cliente gRPC descubra endereços de servidores dinamicamente (ex: consultando DNS, Consul, ou um arquivo de configuração). A Balancer API determina para qual endpoint cada RPC é enviado (round-robin, least-connections, consistent hashing).

**Como funciona em Go:** Implemente `resolver.Builder` e `resolver.Resolver`. O Builder cria instâncias de Resolver. O Resolver envia `resolver.State{Addresses: [...]}` para o gRPC-Go via `cc.UpdateState()`. O Balancer recebe os endereços atualizados e decide o roteamento.

**Por que importa na prática:** Service discovery dinâmico é essencial em Kubernetes onde pods têm IPs efêmeros. Sem resolver customizado, você precisaria de DNS com TTL muito curto ou de um service mesh. Com resolver, você integra com qualquer backend de discovery.

**Armadilha comum:** Resolver que envia updates apenas na criação e não observa mudanças — instâncias que aparecem ou somem não são detectadas. O Resolver deve ter um watcher que atualiza `cc.UpdateState()` quando o conjunto de endereços muda.

**Exercício mental:** Você tem 5 réplicas de um serviço. Duas ficam lentas. Round-robin vai continuar enviando 2/5 do tráfego para elas. Como least-connections resolve isso naturalmente?

---

### Consistent hashing com vnodes

**O que é:** Consistent hashing distribui chaves por um "anel" de hash. Cada nó ocupa uma posição no anel. Uma chave é roteada ao próximo nó no sentido horário. Quando um nó entra ou sai, apenas as chaves entre o nó predecessor e o nó afetado se movem — não todas as chaves.

**Como funciona em Go:** Implemente o anel como slice de hashes ordenados (uint32). Para cada nó físico, crie N vnodes (nós virtuais) com hashes diferentes. Mais vnodes = distribuição mais uniforme. Para resolver uma chave: calcule o hash, faça binary search no slice para encontrar o próximo vnode.

**Por que importa na prática:** Usado para roteamento com afinidade (ex: mesmo usuário sempre vai para o mesmo backend de cache). Vnodes garantem que nós com capacidades diferentes possam ter mais ou menos vnodes proporcionalmente.

**Armadilha comum:** Poucos vnodes por nó físico — sem vnodes, adicionar/remover um nó desequilibra drasticamente a distribuição. Com 100-200 vnodes por nó físico, a distribuição fica estatisticamente uniforme.

**Exercício mental:** Você tem 3 servidores com consistent hashing e 10 vnodes cada. Um servidor sai. Quantas chaves (em porcentagem) precisam ser remapeadas? Compare com hashing simples `hash(key) % N`.

---

### Hedged requests

**O que é:** Estratégia de latência tail: ao enviar uma request a um servidor, se ela não responder em P95 milissegundos, envie a mesma request para uma segunda réplica. Use a resposta que chegar primeiro e cancele a outra. Reduz P99 com custo de ~5% de load extra.

**Como funciona em Go:** Envia a primeira request. Após `hedgeDelay`, envia a segunda com um novo contexto. Usa um channel para receber a primeira resposta que chegar. Cancela as demais via context.

**Por que importa na prática:** P99 frequentemente é causado por um servidor momentaneamente sobrecarregado ou em GC pause. Hedging resolve essa causa sem aumentar infraestrutura — só aumenta tráfego nas caudas, onde justamente há capacidade disponível nas outras réplicas.

**Armadilha comum:** Hedging de operações não-idempotentes (escrita, pagamento) — isso pode causar processamento duplo. Hedging é para operações de leitura ou operações idempotentes apenas.

**Exercício mental:** Seu P50 é 10ms e P99 é 500ms. Com hedging após 20ms, qual o P99 esperado? Por quê o P50 quase não muda com hedging?

---

### Propagação de trace ID via metadata gRPC

**O que é:** Distributed tracing conecta operações em múltiplos serviços num único trace. O trace ID gerado na borda do sistema é propagado para cada serviço downstream via headers HTTP ou metadata gRPC, permitindo correlacionar logs e spans num único request.

**Como funciona em Go:** Interceptors gRPC (UnaryClientInterceptor e UnaryServerInterceptor) injetam e extraem o trace ID do metadata. No servidor: extrai `x-trace-id` do metadata incoming, coloca no context. No cliente: pega do context, injeta no metadata outgoing.

**Por que importa na prática:** Sem trace ID propagado, debugar um erro que cruzou 5 serviços é impossível — você só vê o sintoma no último serviço, não a causa no primeiro. Com trace ID, filtra todos os logs de todos os serviços por um único ID.

**Armadilha comum:** Gerar novo trace ID no primeiro serviço mas não propagá-lo nos serviços internos chamados — você vê o trace no gateway mas não nos microserviços. Todos os clientes gRPC precisam do interceptor de propagação.

**Exercício mental:** Um request passa por 6 serviços. O serviço 4 retorna erro. Sem trace propagation, quais logs você precisa correlacionar manualmente? Com trace propagation, como muda o diagnóstico?

---

### GraphQL com gqlgen: resolvers, dataloaders

**O que é:** GraphQL é uma query language onde o cliente especifica exatamente os dados que precisa. gqlgen é o gerador de código Go que usa o schema GraphQL para gerar interfaces Go type-safe. Resolvers são funções que retornam cada campo. Dataloaders resolvem o N+1 problem.

**Como funciona em Go:** Define o schema em `.graphql`, executa `go run github.com/99designs/gqlgen generate`. gqlgen gera interfaces de resolver. Você implementa cada resolver. Para N+1: use DataLoader que batcha queries individuais numa query única com IN clause.

**Por que importa na prática:** N+1 é o problema mais crítico em GraphQL: um resolver de lista que chama outro resolver para cada item faz 1+N queries. DataLoader resolve batching e caching dentro de um mesmo request.

**Armadilha comum:** DataLoader sem escopo de request — se o DataLoader viver além do request, ele vai retornar dados cacheados de requests anteriores para outros usuários. Crie um novo DataLoader por request, não por servidor.

**Exercício mental:** Uma query GraphQL retorna 100 posts, cada um com o autor. Sem DataLoader, quantas queries ao banco? Com DataLoader que batcha por request, quantas?

---

## Infraestrutura

### Redis: strings, sorted sets, Lua scripts, atomicidade

**O que é:** Redis é um store em memória. Strings são o tipo básico (GET/SET/INCR). Sorted sets associam membros a scores float, mantendo ordenação (ZADD/ZRANGEBYSCORE/ZRANK). Lua scripts executam múltiplas operações atomicamente no servidor Redis.

**Como funciona em Go:** Com `go-redis`: `rdb.Set(ctx, key, value, ttl)`, `rdb.ZAdd(ctx, key, Z{Score, Member})`. Lua via `rdb.Eval(ctx, script, keys, args)`. Lua scripts são atômicos porque Redis é single-threaded — o script roda sem interrupção.

**Por que importa na prática:** Sorted sets são a estrutura ideal para filas de prioridade, leaderboards, e rate limiting com sliding window. Lua resolve o problema de race conditions em operações compostas (check-then-act) sem precisar de WATCH/MULTI/EXEC.

**Armadilha comum:** Lua scripts que fazem chamadas bloqueantes (ex: chamadas de rede) dentro do script — Redis fica bloqueado para todos os outros clientes durante a execução. Lua no Redis deve ser puramente computacional sobre os dados do Redis.

**Exercício mental:** Você quer implementar `GETSET` atômico (pega o valor atual e seta o novo) usando Lua. Qual a diferença de usar `GET` + `SET` sequenciais em vez de Lua em ambiente com múltiplos clientes?

---

### Docker multi-stage build

**O que é:** Multi-stage build usa múltiplos estágios `FROM` num único Dockerfile. O estágio de build tem todas as ferramentas (compilador, dependências). O estágio final copia apenas o binário compilado. A imagem final não contém o toolchain — fica mínima.

**Como funciona em Go:** Estágio 1 (`FROM golang:1.22 AS builder`): copia o código, instala dependências, compila com `CGO_ENABLED=0 GOOS=linux go build`. Estágio 2 (`FROM scratch` ou `FROM alpine`): copia apenas o binário do estágio builder via `COPY --from=builder`.

**Por que importa na prática:** Uma imagem Go sem multi-stage tem ~800MB (imagem base + toolchain). Com multi-stage e `FROM scratch`, chega a ~10MB (apenas o binário estático). Imagens menores = pull mais rápido, menor superfície de ataque, menos CVEs.

**Armadilha comum:** Compilar sem `CGO_ENABLED=0` e usar `FROM scratch` — binários com CGO têm dependências de libc que não existem no scratch. Compile sem CGO para binário totalmente estático.

**Exercício mental:** Seu binário Go precisa fazer chamadas TLS para APIs externas. Usando `FROM scratch`, quais arquivos além do binário você precisa copiar para a imagem final funcionar?

---

### Kubernetes: Deployment, Service, HPA, probes

**O que é:** Deployment gerencia pods com réplicas e rolling updates. Service expõe os pods via IP estável e load balancing. HPA (Horizontal Pod Autoscaler) escala o número de réplicas automaticamente baseado em CPU/memória/métricas customizadas. Probes definem saúde do pod.

**Como funciona em Go:** Seu servidor expõe `/healthz` (liveness: está vivo?) e `/readyz` (readiness: pode receber tráfego?). liveness probe falha → pod é reiniciado. readiness probe falha → pod é removido do Service (para de receber tráfego) sem ser reiniciado. HPA usa a metrics API para escalar.

**Por que importa na prática:** Liveness e readiness separados são críticos. Durante inicialização (carregando cache, conectando ao banco), readiness deve falhar para o pod não receber tráfego antes de estar pronto. Após inicialização, pode falhar liveness se travar completamente.

**Armadilha comum:** Implementar apenas liveness probe e ignorar readiness — um pod que está inicializando mas ainda não pronto recebe tráfego e retorna erros. Implemente ambas com semântica correta.

**Exercício mental:** Seu pod leva 30s para inicializar (carrega modelos ML). Quais os valores corretos de `initialDelaySeconds`, `periodSeconds` e `failureThreshold` para não matar o pod durante inicialização?

---

### Leader election com Redis SETNX

**O que é:** Em um cluster de instâncias idênticas, leader election determina qual instância executa uma tarefa que deve rodar em apenas uma delas (ex: scheduler, cron, cleanup). Com Redis, usa-se `SET key value NX EX ttl` — atômico: só seta se não existir (NX), com expiração automática (EX).

**Como funciona em Go:** Instância tenta `SET leader <instance_id> NX EX 10`. Se retorna OK, é o leader. O leader renova o lock periodicamente (a cada 5s, seta EX=10). Se o leader morre, o lock expira e outra instância assume. Todas as instâncias tentam em loop.

**Por que importa na prática:** Sem leader election, múltiplas instâncias executando um scheduler criam jobs duplicados. Com leader election, você tem high availability (outra instância assume se o leader cair) sem duplicação.

**Armadilha comum:** TTL do lock muito longo — se o leader trava, o lock só libera após o TTL, criando indisponibilidade. TTL deve ser curto (5-30s) com renovação frequente. O período de renovação deve ser substancialmente menor que o TTL (ex: renovar a cada TTL/3).

**Exercício mental:** O leader tem TTL de 10s e renova a cada 3s. O leader fica 7s sem resposta (GC pause longa) mas não morreu. O que acontece com o lock durante esses 7s? Outro nó pode assumir?

---

### Graceful shutdown em sistemas distribuídos

**O que é:** Graceful shutdown é o processo de encerrar um serviço sem perder trabalho em andamento: para de aceitar novas conexões/jobs, termina os jobs em andamento, fecha conexões com recursos (banco, broker), e só então encerra o processo.

**Como funciona em Go:** Escute `os.Signal` (`SIGTERM`, `SIGINT`) via `signal.NotifyContext`. Ao receber o sinal, cancele o context raiz do servidor. Use WaitGroup para aguardar goroutines em andamento. Dê timeout para o shutdown (ex: 30s) — após o qual força o encerramento.

**Por que importa na prática:** Kubernetes envia SIGTERM antes de matar o pod. Sem graceful shutdown, requisições em andamento são interrompidas abruptamente (HTTP 502, transações abortadas). Com graceful shutdown, pods são removidos sem erros visíveis ao usuário.

**Armadilha comum:** Fechar o servidor HTTP imediatamente após receber SIGTERM sem `server.Shutdown(ctx)` — requests em andamento são cortadas. Use `http.Server.Shutdown(ctx)` que espera as requisições em andamento terminarem.

**Exercício mental:** Seu servidor tem 100 requests em andamento quando recebe SIGTERM. O timeout de shutdown é 30s. As requests levam em média 2s. O que acontece? E se algumas requests levarem 45s?

---

### OpenTelemetry: traces, métricas, logs correlacionados

**O que é:** OpenTelemetry é o padrão aberto para instrumentação de observabilidade. Um Trace é o percurso completo de um request pelo sistema, composto de Spans. Métricas são agregações numéricas (counters, histogramas, gauges). Logs correlacionados incluem o trace ID para vincular ao trace.

**Como funciona em Go:** SDK `go.opentelemetry.io/otel`. Crie Tracer com `otel.Tracer("service-name")`. Inicie spans com `tracer.Start(ctx, "operation-name")`. Exporte para Jaeger, Zipkin, ou OTLP. Métricas com `otel/metric`. Cada span pode ter atributos, eventos e status de erro.

**Por que importa na prática:** Com logs apenas, você sabe que algo falhou mas não por quê nem onde. Com traces, você vê o tempo em cada etapa (banco, cache, serviço externo). Com spans correlacionados aos logs, você vai do trace para os logs detalhados daquele span específico.

**Armadilha comum:** Instrumentar apenas o servidor HTTP mas não propagar o context para clientes de banco, Redis e gRPC — os spans filhos não aparecem no trace. Todo I/O que recebe context deve propagar o span.

**Exercício mental:** Um request leva 800ms. Seu trace mostra: handler=800ms, db_query=5ms. Os 795ms restantes estão onde? Como você descobriria o que está consumindo esse tempo?

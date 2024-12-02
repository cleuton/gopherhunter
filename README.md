![](./images/screenshot.png)

# Ajude o Gopher a escapar dos inimigos

[**Cleuton Sampaio**](https://linkedin.com/in/cleutonsampaio)

[**ENGLISH VERSION**](./english.md)

**Gopher** é a mascote da linguagem **Go**. Neste jogo simples, programado em **Go**, demonstro as principais técnicas de games 2D utilizando a biblioteca [**Pixel**](https://github.com/gopxl/pixel/tree/main). 

## Para compilar

Você precisa ter a plataforma **Go** instalada (versão 1.23 ou superior). Este jogo foi criado no S.O. **Ubuntu** 24.04 LTS, e testado com **MacOS** 15 (Sequoia). Não testei com Microsoft **Windows** mas é possível rodar nela. Consulte as instruções da biblioteca [**Pixel**](https://github.com/gopxl/pixel/blob/main/docs/Compilation/Building-Pixel-on-Windows.md).

Clone o repositório, abra a pasta `src` e digite o comando: 

```shell
go mod tidy
``` 

Depois é só compilar ou executar o programa: 

```shell
go run main.go
``` 

## Explicação breve do código fonte

Este jogo, chamado **"Gopher Hunter"**, é um exemplo de um jogo 2D criado com a biblioteca **Pixel**, que é uma biblioteca para gráficos 2D em Go. Vamos começar com uma visão geral do jogo:

### **Estrutura e Conceitos do Jogo**

1. **Componentes Básicos do Jogo 2D:**
   - **Sprites**: Imagens são carregadas e divididas em partes menores (sprites) que representam personagens e elementos do jogo, como o jogador (um gopher), NPCs (serpentes, caranguejos, xícaras) e o cenário.
   - **Cenário de Fundo**: O jogo inclui um fundo que se move para criar a ilusão de movimento enquanto o jogador permanece em posição fixa.

2. **Personagem Principal:**
   - O **jogador** pode realizar ações como pular (`KeyUp`) ou reduzir a velocidade (`KeyLeft`). Essas ações afetam seu comportamento no jogo.
   - O movimento do jogador é baseado em física simples, como gravidade para controlar o salto e limites de altura.

3. **NPCs (Non-Player Characters):**
   - Existem 3 tipos de NPCs: serpentes, caranguejos e xícaras.
   - Cada NPC tem um comportamento único. Por exemplo:
     - Caranguejos (**Rust**) são mais rápidos e podem pular.
     - Serpentes (**Python**) apenas se movem horizontalmente, mas são compridas.
     - Chícaras (**Java**) voam em velocidades diferentes.
   - Todos os NPCs compartilham propriedades comuns, como posição e detecção de colisão.

4. **Colisões:**
   - O jogo verifica colisões entre o jogador e os NPCs. Se uma colisão for detectada, o jogo exibe uma tela de "Game Over" com a opção de reiniciar ou sair.

5. **Lógica de Jogo:**
   - O jogo roda em um loop principal, onde:
     - O tempo (`dt`) é calculado para garantir atualizações suaves.
     - Elementos do cenário e NPCs são movidos e desenhados.
     - Novos NPCs são lançados periodicamente.
     - Entradas do jogador são capturadas para realizar ações.

6. **Recursos Gráficos e Janelas:**
   - A biblioteca Pixel é usada para criar uma janela, desenhar sprites e manipular gráficos 2D.
   - Texto é exibido na tela usando uma fonte básica para instruções e informações como o tempo de jogo.

7. **Funções Importantes:**
   - **`loadPicture`**: Carrega imagens e as transforma em sprites.
   - **`move` e `draw`**: Controlam o movimento e a renderização de NPCs e do jogador.
   - **`collide`**: Verifica colisões usando retângulos delimitadores.

### **Conceitos Associados a Jogos 2D**
- **Scroll Lateral**: O fundo se move para criar a sensação de progresso, um conceito comum em jogos 2D.
- **Física Simples**: Inclui aceleração vertical para o salto do jogador e limites horizontais para elementos.
- **Loop Principal**: O jogo é estruturado em um loop contínuo que processa entradas, atualiza estados e redesenha a tela.
- **Spritesheet**: Divisão de imagens maiores em partes menores para animações e reutilização eficiente de recursos gráficos.

Este jogo demonstra como estruturar um jogo 2D básico, com foco em organização, física simples e interatividade, usando a biblioteca **Pixel** como base gráfica.

## Detalhes de implementação

A implementação do jogo com a biblioteca **Pixel** é bem estruturada para aproveitar seus recursos gráficos e de gerenciamento de eventos. Vamos associar as funcionalidades do jogo com os recursos fornecidos pelo Pixel:

### **1. Gerenciamento de Janelas e Configurações**
- **Pixel fornece:** Criação e controle da janela de jogo por meio do pacote `opengl`.
- **No código:**
  - A janela é configurada e criada com `opengl.WindowConfig` e `opengl.NewWindow`, definindo dimensões, posição, título e sincronização vertical (VSync).
  - Exemplo:
    ```go
    cfg := opengl.WindowConfig{
        Title:    "Gopher Hunter",
        Bounds:   pixel.R(0, 0, 1024, 768),
        Position: pixel.V(posX, posY),
        VSync:    true,
    }
    win, err := opengl.NewWindow(cfg)
    ```
  - **Função associada:** `run()` inicializa e controla o loop principal do jogo, utilizando a janela para desenhar os elementos.

### **2. Gerenciamento de Sprites**
- **Pixel fornece:** O tipo `pixel.Sprite` para representar imagens e métodos para desenhá-las na tela.
- **No código:**
  - Sprites são carregados de imagens (com `loadPicture`) e criados a partir de **spritesheets**, que dividem imagens maiores em pedaços menores.
  - Exemplo:
    ```go
    snakeSprites = append(snakeSprites, pixel.R(x, y, x+128, y+31))
    ```
    Cria subimagens (sprites) das serpentes.
  - Sprites são desenhados na tela com o método `Draw`:
    ```go
    element.Draw(win, matrices[i])
    ```

### **3. Movimento e Transformações**
- **Pixel fornece:** O tipo `pixel.Matrix` para realizar transformações como translação, rotação e escalonamento.
- **No código:**
  - Elementos como NPCs e o cenário usam transformações para se mover horizontalmente ou para simular o movimento da tela.
  - O método `Moved` cria uma matriz de translação:
    ```go
    matrices[i] = matrices[i].Moved(pixel.V(-backSpeedFactor*dt, 0))
    ```
    Move o cenário para a esquerda.

### **4. Manipulação de Eventos**
- **Pixel fornece:** Métodos para detectar entradas do teclado e do mouse.
- **No código:**
  - Os eventos capturam ações do jogador, como pular (`KeyUp`) ou reduzir velocidade (`KeyLeft`):
    ```go
    if win.JustPressed(pixel.KeyUp) {
        player.isJumping = true
    }
    ```
  - O método `JustPressed` verifica se uma tecla foi pressionada.

### **5. Detecção de Colisões**
- **Pixel fornece:** Ferramentas para manipular vetores (`pixel.V`) e retângulos (`pixel.Rect`), essenciais para verificar colisões.
- **No código:**
  - Colisões são verificadas pelo método `Intersect`, que calcula a interseção entre dois retângulos:
    ```go
    collision := elementRect.Intersect(rect)
    return collision.Area() > 0
    ```

### **6. Controle de Tempo**
- **Pixel fornece:** Integração com pacotes padrão de Go (`time`) para medir e controlar o tempo.
- **No código:**
  - O intervalo de tempo (`dt`) entre quadros é calculado para ajustar o movimento suavemente:
    ```go
    dt := time.Since(last).Seconds()
    ```
  - Usado para ajustar a velocidade dos NPCs e a duração de animações.

### **7. Texto na Tela**
- **Pixel fornece:** O pacote `text` para criar e manipular textos.
- **No código:**
  - Textos como tempo decorrido ou mensagens de fim de jogo são desenhados com `text.New`:
    ```go
    fmt.Fprintf(seconds, secondsText, secondsRunning)
    seconds.Draw(win, pixel.IM.Scaled(seconds.Orig, 2))
    ```

### **8. Estruturas e Funções Principais**
- **Estruturas**:
  - `Player` e NPCs (`Snake`, `Crab`, `Cup`) compartilham uma estrutura base (`CommonNpcProperties`) com propriedades comuns como posição, tamanho e velocidade.
  - Cada tipo de NPC implementa o comportamento de movimento (`move`), desenho (`draw`) e colisão (`collide`).
- **Funções**:
  - **`move(dt float64)`**: Atualiza a posição dos NPCs e do jogador com base no tempo.
  - **`draw(pixel.Target)`**: Renderiza o elemento na tela.
  - **`loadPicture(path string)`**: Carrega imagens do disco para serem usadas como sprites.

### **Resumindo o Papel do Pixel na Implementação**
A biblioteca **Pixel** simplifica o desenvolvimento do jogo, fornecendo ferramentas para:
1. Criar e gerenciar a janela.
2. Trabalhar com sprites e animações.
3. Detectar eventos do teclado.
4. Implementar colisões e movimentação suave.
5. Adicionar elementos como texto e gráficos dinâmicos.

Essa combinação de recursos permite que o jogo seja estruturado em torno de um **loop principal**, com foco na interação do jogador, animações fluidas e um ambiente visual rico.

## Estrutura do código

A estrutura do código segue uma abordagem típica para jogos 2D em que os elementos são bem modularizados, com responsabilidades claras e reutilização de componentes. Aqui está a explicação da estrutura:

### **1. Configuração Inicial e Variáveis Globais**
- **Propósito:** Definir configurações e armazenar o estado global do jogo.
- **Elementos Principais:**
  - Variáveis globais, como `player`, `npcs`, `elements`, e `secondsRunning`, armazenam o estado do jogador, inimigos, elementos visuais e tempo de jogo.
  - Configurações de velocidade, posição inicial e propriedades dos elementos são declaradas para fácil ajuste.
  - Exemplo:
    ```go
    backSpeedFactor := 50.0
    crabSpeed := 120.0
    playerJumpLimit := 500.0
    ```

### **2. Estruturas de Dados**
- **Propósito:** Modelar os elementos principais do jogo.
- **Elementos Principais:**
  - **`CommonNpcProperties`:** Uma estrutura base que contém propriedades comuns para NPCs e o jogador, como posição, velocidade e tamanho.
    ```go
    type CommonNpcProperties struct {
        sprite1       *pixel.Sprite
        sprite2       *pixel.Sprite
        position      pixel.Vec
        height        float64
        width         float64
        speed         float64
        horizontalWay float64
        inverted      bool
    }
    ```
  - **NPCs Específicos:** Estruturas como `Crab`, `Snake` e `Cup` estendem `CommonNpcProperties` para implementar comportamentos únicos, como pular no caso do `Crab`.
  - **`Player`:** O jogador tem comportamentos específicos, como pular e reduzir velocidade, com sua própria lógica de movimento.

### **3. Interfaces**
- **Propósito:** Garantir que todos os NPCs compartilhem um conjunto de comportamentos básicos.
- **Elementos Principais:**
  - A interface `Npc` define métodos que cada NPC deve implementar:
    ```go
    type Npc interface {
        move(dt float64) bool
        draw(pixel.Target)
        collide(pixel.Rect) bool
    }
    ```
  - Isso facilita o polimorfismo e a manipulação genérica de NPCs no jogo.

### **4. Funções de Comportamento**
- **Propósito:** Implementar a lógica principal dos NPCs e do jogador.
- **Elementos Principais:**
  - **`move(dt float64) bool`:** Atualiza a posição dos elementos com base no tempo.
  - **`draw(pixel.Target)`:** Desenha o elemento na tela.
  - **`collide(pixel.Rect) bool`:** Verifica colisões entre elementos.
  - Exemplo de lógica de colisão:
    ```go
    func (c CommonNpcProperties) collide(rect pixel.Rect) bool {
        lowerLeft := pixel.V(c.position.X-c.width/2, c.position.Y-c.height/2)
        upperRight := pixel.V(c.position.X+c.width/2, c.position.Y+c.height/2)
        elementRect := pixel.R(lowerLeft.X, lowerLeft.Y, upperRight.X, upperRight.Y)
        collision := elementRect.Intersect(rect)
        return collision.Area() > 0
    }
    ```

### **5. Funções Auxiliares**
- **Propósito:** Realizar tarefas específicas, como carregar imagens ou inicializar elementos.
- **Elementos Principais:**
  - **`loadPicture(path string)`**: Carrega imagens e as converte em `pixel.Picture` para uso no jogo.
  - **NPC Constructors (`NewSnake`, `NewCrab`, `NewCup`)**:
    - Criam instâncias de NPCs, carregando spritesheets e definindo propriedades específicas.
    - Exemplo:
      ```go
      func NewSnake() *Snake {
          return &Snake{
              CommonNpcProperties{
                  sprite1: pixel.NewSprite(snakeSpriteSheet, snakeSprites[0]),
                  sprite2: pixel.NewSprite(snakeSpriteSheet, snakeSprites[1]),
                  position: pixel.V(1024, 200+31/2),
                  speed: snakeSpeed,
              },
          }
      }
      ```

### **6. Inicialização do Jogo**
- **Propósito:** Configurar o ambiente inicial antes do loop principal.
- **Elementos Principais:**
  - A função `initGame()` reseta variáveis globais e reconfigura o estado do jogo.
  - Exemplo:
    ```go
    func initGame() {
        elements = []*pixel.Sprite{}
        npcs = []Npc{}
        secondsRunning = 0.0
    }
    ```

### **7. Loop Principal do Jogo**
- **Propósito:** Controlar o fluxo do jogo, incluindo entrada do jogador, atualizações e renderização.
- **Elementos Principais:**
  - A função `run()` implementa o loop principal, que executa até que a janela seja fechada:
    - Calcula o tempo (`dt`) entre quadros para movimentos suaves.
    - Processa entradas do teclado.
    - Atualiza posições dos elementos (jogador, NPCs e cenário).
    - Verifica colisões e termina o jogo em caso de impacto.
    - Desenha todos os elementos na tela.
  - Exemplo:
    ```go
    for !win.Closed() {
        dt := time.Since(last).Seconds()
        last = time.Now()
        player.move(dt)
        player.draw(win)
    }
    ```

### **8. Fim de Jogo e Reinício**
- **Propósito:** Exibir uma tela de "Game Over" e permitir que o jogador reinicie ou saia.
- **Elementos Principais:**
  - Após uma colisão, o jogo exibe uma mensagem e opções para continuar:
    ```go
    if npc.collide(player.rect()) {
        for !win.Closed() {
            basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 3))
            if win.JustPressed(pixel.KeyY) {
                initGame()
            }
        }
    }
    ```

### **Resumo da Estrutura**
- **Organização modular:** Com separação clara entre inicialização, comportamento de NPCs, loop principal e gerenciamento de estados.
- **Reutilização:** Uso extensivo de estruturas comuns e interfaces para evitar repetição de código.
- **Ciclo de jogo bem definido:** Inicialização → Loop Principal → Fim/Reinício.

Essa estrutura modular facilita a manutenção e extensão do jogo, permitindo a adição de novos elementos ou mecânicas sem grandes alterações.
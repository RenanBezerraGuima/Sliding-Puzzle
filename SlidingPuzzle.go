package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Puzzle struct {
	matriz      [3][3]int
	vazioLinha  int
	vazioColuna int
}

type No struct {
	p            *Puzzle
	pai          *No
	profundidade int
	movimento    string
}

var estadoObjetivo = [3][3]int{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 0},
}

func (p *Puzzle) inicializar() {
	p.matriz = estadoObjetivo
	p.vazioLinha = 2
	p.vazioColuna = 2
}

func (p *Puzzle) obterMovimentosPossiveis() []string {
	var movimentos []string

	if p.vazioLinha > 0 {
		movimentos = append(movimentos, "cima")
	}

	if p.vazioLinha < 2 {
		movimentos = append(movimentos, "baixo")
	}

	if p.vazioColuna > 0 {
		movimentos = append(movimentos, "esquerda")
	}

	if p.vazioColuna < 2 {
		movimentos = append(movimentos, "direita")
	}

	return movimentos
}

func (p *Puzzle) executarMovimento(movimento string) bool {
	novaLinha, novaColuna := p.vazioLinha, p.vazioColuna

	switch movimento {
	case "cima":
		novaLinha--
	case "baixo":
		novaLinha++
	case "esquerda":
		novaColuna--
	case "direita":
		novaColuna++
	default:
		return false
	}

	if novaLinha < 0 || novaLinha >= 3 || novaColuna < 0 || novaColuna >= 3 {
		return false
	}

	p.matriz[p.vazioLinha][p.vazioColuna] = p.matriz[novaLinha][novaColuna]
	p.matriz[novaLinha][novaColuna] = 0
	p.vazioColuna = novaColuna
	p.vazioLinha = novaLinha

	return true
}

// Começa com o estado objetivo e randomiza movimentos
// para garantir que o problema seja solúvel
func (p *Puzzle) randomizar(movimentos int) {
	p.inicializar()

	for i := 0; i < movimentos; i++ {
		movimentosPossiveis := p.obterMovimentosPossiveis()
		if len(movimentosPossiveis) > 0 {
			movimento := movimentosPossiveis[rand.Intn(len(movimentosPossiveis))]
			p.executarMovimento(movimento)
		}
	}
}

func (p *Puzzle) estaCompleto() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if p.matriz[i][j] != estadoObjetivo[i][j] {
				return false
			}
		}
	}
	return true
}

func (p *Puzzle) obterChave() string {
	chave := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			chave += fmt.Sprintf("%d", p.matriz[i][j]) // ?
		}
	}
	return chave
}

func (p *Puzzle) copiar() *Puzzle {
	novoPuzzle := &Puzzle{
		vazioLinha:  p.vazioLinha,
		vazioColuna: p.vazioColuna,
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			novoPuzzle.matriz[i][j] = p.matriz[i][j]
		}
	}
	return novoPuzzle
}

func resolverBuscaEmLargura(puzzleInicial *Puzzle) (*No, int) {
	if puzzleInicial.estaCompleto() {
		return &No{p: puzzleInicial, pai: nil, profundidade: 0, movimento: ""}, 0
	}

	fila := []*No{{p: puzzleInicial, pai: nil, profundidade: 0, movimento: "inicial"}}
	visitados := make(map[string]bool)
	visitados[puzzleInicial.obterChave()] = true
	nosExpandidos := 0

	for len(fila) > 0 {
		noAtual := fila[0]
		fila = fila[1:]
		nosExpandidos++

		movimentos := noAtual.p.obterMovimentosPossiveis()
		for _, movimento := range movimentos {
			novoPuzzle := noAtual.p.copiar()
			novoPuzzle.executarMovimento(movimento)
			chave := novoPuzzle.obterChave()

			if !visitados[chave] {
				visitados[chave] = true
				novoNo := &No{
					p:            novoPuzzle,
					pai:          noAtual,
					profundidade: noAtual.profundidade + 1,
					movimento:    movimento,
				}

				if novoPuzzle.estaCompleto() {
					return novoNo, nosExpandidos
				}

				fila = append(fila, novoNo)
			}
		}
	}
	return nil, nosExpandidos // Não encontrou solução
}

func (p *Puzzle) calcularQtdNosForaDePosicao() int {
	qtd := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if p.matriz[i][j] != estadoObjetivo[i][j] && p.matriz[i][j] != 0 {
				qtd++
			}
		}
	}
	return qtd
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p *Puzzle) calcularDistanciaManhattan() int {
	distancia := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			valor := p.matriz[i][j]
			if valor != 0 {
				linhaObjetivo := (valor - 1) / 3
				colunaObjetivo := (valor - 1) % 3

				distancia += abs(i-linhaObjetivo) + abs(j-colunaObjetivo)
			}
		}
	}
	return distancia
}

func (p *Puzzle) calcularHeuristica(tipoHeuristica string) int {
	switch tipoHeuristica {
	case "manhattan", "Manhattan", "H2", "h2":
		return p.calcularDistanciaManhattan()
	case "fora_posicao", "foraPosicao", "ForaPosicao", "H1", "h1":
		return p.calcularQtdNosForaDePosicao()
	default:
		fmt.Printf("Heurística '%s' não reconhecida, utilizando Manhattan como padrão. \n", tipoHeuristica)
		return p.calcularDistanciaManhattan()
	}
}

func resolverAEstrela(puzzleInicial *Puzzle, heuristica string) (*No, int) {
	if puzzleInicial.estaCompleto() {
		return &No{p: puzzleInicial, pai: nil, profundidade: 0, movimento: ""}, 0
	}

	listaAberta := []*No{{p: puzzleInicial, pai: nil, profundidade: 0, movimento: "inicial"}}
	visitados := make(map[string]bool)
	nosExpandidos := 0

	for len(listaAberta) > 0 {
		indiceMinimo := 0
		menorF := listaAberta[0].profundidade + listaAberta[0].p.calcularHeuristica(heuristica)

		for i := 1; i < len(listaAberta); i++ {
			f := listaAberta[i].profundidade + listaAberta[i].p.calcularHeuristica(heuristica)
			if f < menorF {
				menorF = f
				indiceMinimo = i
			}
		}

		noAtual := listaAberta[indiceMinimo]
		listaAberta = append(listaAberta[:indiceMinimo], listaAberta[indiceMinimo+1:]...)
		chaveAtual := noAtual.p.obterChave()

		if visitados[chaveAtual] {
			continue
		}
		visitados[chaveAtual] = true
		nosExpandidos++

		if noAtual.p.estaCompleto() {
			return noAtual, nosExpandidos
		}

		movimentos := noAtual.p.obterMovimentosPossiveis()
		for _, movimento := range movimentos {
			novoPuzzle := noAtual.p.copiar()
			novoPuzzle.executarMovimento(movimento)
			chave := novoPuzzle.obterChave()

			if !visitados[chave] {
				novoNo := &No{
					p:            novoPuzzle,
					pai:          noAtual,
					profundidade: noAtual.profundidade + 1,
					movimento:    movimento,
				}
				listaAberta = append(listaAberta, novoNo)
			}
		}
	}
	return nil, nosExpandidos
}

func (p *Puzzle) imprimir() {
	fmt.Println("┌───┬───┬───┐")
	for i := 0; i < 3; i++ {
		fmt.Print("│")
		for j := 0; j < 3; j++ {
			if p.matriz[i][j] == 0 {
				fmt.Print("   │")
			} else {
				fmt.Printf(" %d │", p.matriz[i][j])
			}
		}
		if i == 2 {
			fmt.Println()
			fmt.Println("└───┴───┴───┘")
		} else {
			fmt.Println()
			fmt.Println("├───┼───┼───┤")
		}
	}
}

func reconstruirCaminho(no *No) []string {
	var caminho []string
	atual := no

	for atual != nil && atual.pai != nil {
		caminho = append([]string{atual.movimento}, caminho...)
		atual = atual.pai
	}

	return caminho
}

func executarSolucao(puzzle *Puzzle, caminho []string) {
	fmt.Println("Executando solução passo a passo:")
	fmt.Println("\nEstado inicial:")
	puzzle.imprimir()

	for i, movimento := range caminho {
		fmt.Printf("\nPasso %d - Movimento: %s\n", i+1, movimento)
		puzzle.executarMovimento(movimento)
		puzzle.imprimir()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	var p Puzzle
	var ph1 Puzzle
	var ph2 Puzzle
	fmt.Println("Estado Objetivo:")
	p.inicializar()
	p.imprimir()
	movimentos := 200
	fmt.Printf("Estado Inicial randomizado com %d movimentos aleátorios\ndo estado objetivo, garantido que este seja solúvel:\n", movimentos)
	p.randomizar(movimentos)
	p.imprimir()

	fmt.Printf("\nHeurística 1 (Quantidade de números fora de posição): %d\n", p.calcularQtdNosForaDePosicao())
	fmt.Printf("Heurística 2 (Distância Manhattan dos números da sua posição objetivo): %d\n", p.calcularDistanciaManhattan())

	ph1 = *p.copiar()
	ph2 = *p.copiar()

	noObjetivoBEL, nosExpandidosBEL := resolverBuscaEmLargura(&p)
	// executarSolucao(&p, reconstruirCaminho(noObjetivoBEL))

	noObjetivoh1, nosExpandidosh1 := resolverAEstrela(&ph1, "h1")

	noObjetivoh2, nosExpandidosh2 := resolverAEstrela(&ph2, "h2")

	fmt.Printf("Quantidade de nós expandidos pela busca em Largura: %d\n", nosExpandidosBEL)
	fmt.Printf("Profundidade do nó (quantidade de movimentos) Busca Em Largura (Solução ótima): %d\n", noObjetivoBEL.profundidade)
	fmt.Printf("Quantidade de nós expandidos pelo A* h1: %d\n", nosExpandidosh1)
	fmt.Printf("Profundidade do nó (quantidade de movimentos) A* h1: %d\n", noObjetivoh1.profundidade)
	fmt.Printf("Quantidade de nós expandidos pelo A* h2: %d\n", nosExpandidosh2)
	fmt.Printf("Profundidade do nó (quantidade de movimentos): A* h2 %d\n", noObjetivoh2.profundidade)
}

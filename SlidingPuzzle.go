package main

import (
	"fmt"
	"math/rand"
)

type Puzzle struct {
	matriz [3][3]int
}

func (p *Puzzle) randomizar() {
	for i := 1; i < 9; i++ {
		for {
			random1 := rand.Intn(3)
			random2 := rand.Intn(3)
			if p.matriz[random1][random2] == 0 {
				p.matriz[random1][random2] = i
				break
			}
		}

	}
}

func (p *Puzzle) imprimirMatriz() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Print(p.matriz[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	var p Puzzle
	p.imprimirMatriz()
	p.randomizar()
	p.imprimirMatriz()
}

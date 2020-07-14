package main

import (
	"fmt"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

func TestUnSat() { 	// Creates a test of three literals, and then three literals to prove that the solve is unsat.
	g := gini.New()
	g.Add(z.Lit(3))
	g.Add(0)
	g.Add(z.Lit(3).Not())
	g.Add(0)
	if g.Solve() != -1 {
		fmt.Printf("basic add unsat failed.\n")
	} else {
		fmt.Printf("The Solve is UNSAT!\n")
	}
}
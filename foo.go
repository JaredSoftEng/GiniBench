package main

import (
	"fmt"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
	"strconv"
)

func main() {
	testUnSat()
	testSat()
}

func testUnSat() { 	// Creates a test of three literals, and then three literals to prove that the solve is unsat.
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

func testSat() {
	g := gini.New()
	g.Add(z.Lit(1))
	g.Add(0)
	g.Add(z.Lit(2).Not())
	g.Add(0)
	if g.Solve() == 1 {
		fmt.Printf("The solution is SAT: ")
		fmt.Printf("Lit 1: %s, Lit 2: %s\n", strconv.FormatBool(g.Value(1)), strconv.FormatBool(g.Value(2)))
	}
}

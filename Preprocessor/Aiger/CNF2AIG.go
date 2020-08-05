package CNF2AIG

import (
	"github.com/jaredsofteng/gini"
	"github.com/jaredsofteng/gini/logic/aiger"
	"io"
	"log"
	"os"
)

// CNF2AIG takes as an input a CNF problem and attempts to restructure it using the and-inverter graph transformations
// It provides as an output a structure in a logical gate format which corresponds to the input CNF.
// TODO: What do we do with dangling CNF instances
func CNF2AIG(file string) *aiger.T {

}

func Aig2Solve(a aiger.T) *gini.Gini {
	g := gini.New()
	a.C.ToCnf(g)
	return g
}

func ReadAigerBinary(filepath string) *aiger.T {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var r io.Reader
	r = f

	s, err := aiger.ReadBinary(r)
	if err != nil {panic(err)}
	return s
}

func ReadAigerAscii(filepath string) *aiger.T {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var r io.Reader
	r = f

	s, err := aiger.ReadAscii(r)
	if err != nil {panic(err)}
	return s
}
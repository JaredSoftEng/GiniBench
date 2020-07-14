package main

import (
	"fmt"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/bench"
	"github.com/irifrance/gini/z"
	"io"
	"os"
	"path"
	"strconv"
)

func main() {
	s := initializeBench()
	gp := os.Getenv("GOPATH")
	runDir := path.Join(gp, "src/GiniBench/Example_suite/runs")
	run1 := path.Join(runDir, "example1")
	if !bench.IsRunDir(run1) {
		r, _ := bench.NewRun(s, "example1", path.Join(gp, "bin/gini.exe") + " " + path.Join(runDir, s.Insts[0]) + " " + path.Join(runDir, s.Insts[1]) + " " + path.Join(runDir, s.Insts[2]),1000, 1000)
		var index int
		for range s.Insts {
			r.Do(index)
			index += 1
		}
	} else {
		r, _ := bench.OpenRun(s, run1)
		fmt.Print(r.InstTimeout)
	}
}

func initializeBench() *bench.Suite {
	gp := os.Getenv("GOPATH")
	rootDir := path.Join(gp, "src/GiniBench/Example_suite")
	var suite *bench.Suite
	var err error
	if !bench.IsSuiteDir(rootDir) {
		var cnfFile [3]string
		cnfFile[0] = path.Join(gp, "src/GiniBench/Benchmark Problems/Example CNF/example1.cnf")
		cnfFile[1] = path.Join(gp, "src/GiniBench/Benchmark Problems/Example CNF/example2.cnf")
		cnfFile[2] = path.Join(gp, "src/GiniBench/Benchmark Problems/Example CNF/example3.cnf")
		var RunArray = []string{cnfFile[0], cnfFile[1], cnfFile[2]}
		suite, err = bench.CreateSuite(rootDir, RunArray)
	} else {
		suite, err = bench.OpenSuite(rootDir)
	}
	if err != nil {fmt.Print(err.Error())}
	return suite
}

func doRun(suite *bench.Suite) {
	run1, err := bench.NewRun(suite, "Run2", "go run testUnSat.go", 10000, 10000)
	if err != nil {
		fmt.Print(err.Error())
	}
	if run1 != nil {
		fmt.Printf("initializeBench has executed the following: %s", run1.Cmd)
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

func testImportCNF() {
	gp := os.Getenv("GOPATH")
	f, err := os.Open(path.Join(gp, "src/GiniBench/Example_suite/bench-00.cnf"))
	if err != nil {panic(err)}
	defer f.Close()
	var r io.Reader
	r = f
	g, _ := gini.NewDimacs(r)
	if g.Solve() == 1 {
		fmt.Print(g.MaxVar())
	}
}

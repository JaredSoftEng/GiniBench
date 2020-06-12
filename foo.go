package main

import (
	"fmt"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/bench"
	"github.com/irifrance/gini/z"
	"os"
	"path"
	"strconv"
)

func main() {
	s := initializeBench()
	//doRun(s)
	r, err := s.Run("Run3", "go run testUnSat.go", 10000, 10000 )
	if r != nil {
		fmt.Print(r.)
	} else {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	//testSat()

}

func initializeBench() *bench.Suite {
	gp := os.Getenv("GOPATH")
	rootDir := path.Join(gp, "src/GiniBench/suite")
	//rootRunDir := path.Join(rootDir, "runs/Run1")
	//var RunArray []string
	//suite, err := bench.CreateSuite(rootDir, RunArray)
	checkSuite(rootDir)
	//checkRun(rootRunDir)
	suite, err := bench.OpenSuite(rootDir)
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

func checkSuite(path string) {
	check := bench.IsSuiteDir(path)
	if ! check {
		fmt.Printf("The root directory %s is not a valid suite (needs maps, hash and runs subdirectories.", path)
		os.Exit(1)
	}
}

func checkRun(path string) {
	check := bench.IsRunDir(path)
	if ! check {
		fmt.Printf("The 'runs' directory %s is invalid.", path)
		os.Exit(1)
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

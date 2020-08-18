package main

import (
	"GiniBench/Tools"
	"flag"
	"fmt"
	"github.com/jaredsofteng/gini"
	"github.com/jaredsofteng/gini/z"
	os "os"
	"path/filepath"
	"strconv"
	"time"
)

var showUI = flag.Bool("ui", false, "if in a GUI based OS, shows file-picker dialog and opens a CSV of results (must be the only option in the command)")
var timeout = flag.Duration("timeout", time.Second*30, "solver timeout")
var model = flag.Bool("model", false, "if true, prints out the model (default false)")
var nosub = flag.Bool("nosub", false, "skips subsumption (default false)")
var noself = flag.Bool("noself", false, "skips self-subsumption (default false)")
var fullsub = flag.Bool("fullsub", false, "performs full instead of binary subsumption")
var cnf = flag.Bool("cnf", false, "will output a cnf file and not perform a solve")

func main() {
	var (
		help    bool
	)
	flag.BoolVar(&help, "help", false, "displays help")
	flag.Usage = func() {
		p := os.Args[0]
		_, p = filepath.Split(p)
		fmt.Fprintf(os.Stderr, usage, p, p, p, p, p, p, p)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse()
	if flag.NArg() == 0 || flag.NArg() > 1 || *showUI {
		if flag.NArg() == 0 && *showUI {
			fmt.Println("Opening filepicker dialog. -model option will be ignored")
			OpenUI()
			os.Exit(0)
		}
		if flag.NArg() != 0 && *showUI {
			fmt.Printf("The UI option does not accept file references as it exposes a file-picker\n")
			fmt.Printf("Syntax : %s [options] \n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(0)
		}
		if flag.NArg() > 1 && !help {
			fmt.Printf("GiniPre requires a single DIMACS or AIGER input file to perform optimized preprocessing on. For advanced file selection, use the -ui option without any file specified.\n")
			fmt.Fprintf(os.Stderr, "Syntax : %s [options] (*.cnf|*.bz2|*.gz|*.aig|*.aag)\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
	if help {
		fmt.Printf("This is GiniPre version 1.0, a SAT pre-processor by Michael Behr and Jared Lenos.\n")
		fmt.Printf("Syntax : %s [options] (*.cnf|*.bz2|*.gz|*.aig|*.aag)\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *cnf {
		g := readFile(flag.Args()[0])
		if !*noself {
			preprocessSelf(g)
		}
		if !*nosub {
			preprocessSub(g)
		}
		file, _ := os.Create(filepath.Dir(flag.Args()[0]) + filepath.Base(os.Args[0][:len(flag.Args()[0])-len(filepath.Ext(flag.Args()[0]))]) + "-ginipre.cnf")
		g.Write(file)
		os.Exit(0)
	}
	fmt.Printf("c solving %s\n", flag.Args()[0])
	g := readFile(flag.Args()[0])
	if !*noself {
		preprocessSelf(g)
	}
	if !*nosub {
		preprocessSub(g)
	}
	startMem := Tools.TotalMemUsageMB()
	doSolve := g.GoSolve()
	startSolve := time.Now()
	r := doSolve.Try(*timeout)
	endSolve := time.Since(startSolve)
	memConsumed := Tools.TotalMemUsageMB() - startMem
	fmt.Println("c Memory Consumed = " + strconv.FormatUint(memConsumed, 10) + "MB")
	fmt.Println("c Solve Time = " + endSolve.String())
	fmt.Println("c " + printResult(r))
	if *model {
		outputModel(g.MaxVar(), gModel{g})
	}
}

type gModel struct {
	g *gini.Gini
}

func (g gModel) value(m z.Lit) bool {
	return g.g.Value(z.Lit(m))
}

type values interface {
	value(m z.Lit) bool
}

func outputModel(v z.Var, m values) {
	var col = 2
	fmt.Printf("v ")
	for i := z.Var(1); i <= v; i++ {
		n := 0
		for j := i; j > 0; j = j / 10 {
			n++
		}
		t := m.value(i.Pos())
		if !t {
			n++
		}
		if col+n > 78 {
			fmt.Printf("\nv")
			col = 2
		}
		if t {
			fmt.Printf(" %s", i.Pos())
		} else {
			fmt.Printf(" %s", i.Neg())
		}
		col++
		col += n
	}
	fmt.Printf(" 0\n")
}

package main

import (
	"GiniBench/Tools"
	"github.com/irifrance/gini"
	"log"
	"strconv"

	"io"
	"fmt"
	"os"
	"path"
	"time"
)

var logFile string

func main() {
	gp := os.Getenv("GOPATH")
	file1 := path.Join(gp, "/src/Ginibench/Benchmark Problems/Bounded Model Checking/bmc-ibm-1.cnf")
	g := readFile(file1)
	r := solveFile(g)
	var result string
	switch r {
	case 0:
		result = "UNKNOWN"
	case 1:
		result = "SAT"
	case -1:
		result = "UNSAT"
	default:
		result = "ERR IN SOLVE"
	}
	logToFile("Solve result = " + result)
}

func readFile(path string) *gini.Gini {
	startTime := time.Now()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r io.Reader
	r = f
	g, err := gini.NewDimacs(r)
	if err != nil {panic(err)}
	fileReadTime := time.Since(startTime)
	setLogFile(path)
	logToFile("DIMACS parsing time = " + fileReadTime.String())
	return g
}

func solveFile(g *gini.Gini) int {
	startSolve := time.Now()
	startMem := Tools.TotalMemUsageMB()
	startCPU := Tools.CpuUsagePercent()
	doSolve := g.GoSolve()
	result := doSolve.Try(time.Second*30)
	endSolve := time.Since(startSolve)
	logToFile("Solve Time = " + endSolve.String())
	cpuPercent := Tools.CpuUsagePercent()
	logToFile("CPU Usage % = " + strconv.FormatFloat(cpuPercent - startCPU, 'f', 6, 64))
	memConsumed := Tools.TotalMemUsageMB() - startMem
	logToFile("Memory Usage Total = " + string(memConsumed))
	return result
}

func logToFile(s string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { log.Fatal(err)}
	if _, err := fmt.Fprintln(f, s); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func setLogFile(fileIn string) {
	var filename = fileIn
	var filenameExt = path.Ext(filename)
	var newFilenameExt = ".txt"
	newFile := filename[0:len(filename)-len(filenameExt)]
	logFile = newFile + "-result" + newFilenameExt
}
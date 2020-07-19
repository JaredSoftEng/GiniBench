package main

import (
	"GiniBench/Tools"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"github.com/irifrance/gini"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var logFile string

func main() {
	//gp := os.Getenv("GOPATH")
	//file1 := path.Join(gp, "/src/Ginibench/Benchmark Problems/Bounded Model Checking/bmc-ibm-10.cnf")
	file1 := "C:/Users/Jared/Downloads/Main/final/Zhou/queen8-8-9.cnf.bz2"
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
	logToFile(" ") // Space before next run if multiple are done
}

func readFile(filepath string) *gini.Gini {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := readFileReader(f)
	startTime := time.Now()
	g, err := gini.NewDimacs(r)
	if err != nil {panic(err)}
	fileReadTime := time.Since(startTime)
	setLogFile(filepath)
	logToFile("DIMACS parsing time = " + fileReadTime.String())
	return g
}

func readFileReader(f *os.File) io.Reader {
	var r io.Reader
	var e error
	switch path.Ext(f.Name()) {
	case ".gz":
		r, e = gzip.NewReader(f)
		if e != nil {
			log.Fatal(err)
		}
	case ".bz2":
		r = bzip2.NewReader(f)
	case ".cnf":
		r = f
	default:
		log.Fatal("Invalid File format - must be .gz .bz2 or .cnf")
	}
	return r
}


func solveFile(g *gini.Gini) int {
	startMem := Tools.TotalMemUsageMB()
	startCPU := Tools.CpuUsagePercent()
	doSolve := g.GoSolve()
	startSolve := time.Now()
	result := doSolve.Try(time.Second*30)
	endSolve := time.Since(startSolve)
	logToFile("Solve Time = " + endSolve.String())
	cpuPercent := Tools.CpuUsagePercent()
	logToFile("CPU Usage % = " + strconv.FormatFloat(cpuPercent - startCPU, 'f', 6, 64))
	memConsumed := Tools.TotalMemUsageMB() - startMem
	logToFile("Memory Usage Total = " + strconv.FormatUint(memConsumed, 10) + "MB")
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
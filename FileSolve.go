package main

import (
	"GiniBench/Tools"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"github.com/irifrance/gini"
	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var logFile string

func main() {
	startFolder := path.Join(os.Getenv("USERPROFILE"), "Downloads")
	file1, err := dialog.File().SetStartDir(startFolder).Filter("DIMACS File", "bz2", "cnf", "gz").Load()
	if err != nil {
		log.Fatal("File Selection Incomplete")
		return
	}
	ok := dialog.Message("%s", "Import entire Directory?").Title("Import Scope").YesNo()
	if ok {
		files, err := Tools.WalkMatch(path.Dir(file1), "*.cnf")
		bzFiles, err := Tools.WalkMatch(path.Dir(file1), "*.bz")
		gzFiles, err := Tools.WalkMatch(path.Dir(file1), "*.gz")
		files = append(files, bzFiles...)
		files = append(files, gzFiles...)
		if err != nil {
			log.Fatal(err)
		}
		setLogDir(file1)
		for _, f := range files {
			solveMainRoutine(f)
		}
	} else {
		setLogFile(file1)
		solveMainRoutine(file1)
	}
}

func solveMainRoutine(filepath string) {
	logToFile(path.Base(filepath))
	g := readFile(filepath)
	maxSolveTime := time.Second * 30
	r := solveFile(g, maxSolveTime)
	printResult(r)
	_ = open.Start(logFile)
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
			log.Fatal(e.Error())
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


func solveFile(g *gini.Gini, timeout time.Duration) int {
	startMem := Tools.TotalMemUsageMB()
	doSolve := g.GoSolve()
	Tools.CpuUsagePercent(100 * time.Microsecond) // Tracks CPU percent for the next 100 microseconds
	startSolve := time.Now()
	result := doSolve.Try(timeout)
	endSolve := time.Since(startSolve)
	logToFile("Solve Time = " + endSolve.String())
	cpuPercentChange := Tools.CpuUsagePercent(0) // returns difference from last cpu check
	logToFile("CPU Usage % = " + strconv.FormatFloat(cpuPercentChange, 'f', 6, 64))
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

func setLogDir(fileIn string) {
	newFilenameExt := ".txt"
	logFile = path.Dir(fileIn) + "Directory-CNF-Results" + newFilenameExt
}


func printResult(result int) {
	var resultStr string
	switch result {
	case 0:
		resultStr = "UNKNOWN"
	case 1:
		resultStr = "SAT"
	case -1:
		resultStr = "UNSAT"
	default:
		resultStr = "ERR IN SOLVE"
	}
	logToFile("Solve result = " + resultStr)
	logToFile(" ") // Space before next run if multiple are done

}


package main

import (
	"C"
	"GiniBench/Preprocessor/Preprocessor"
	"GiniBench/Preprocessor/pregini"
	"GiniBench/Tools"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"github.com/jaredsofteng/gini"
	"github.com/jaredsofteng/gini/gen"
	"github.com/jaredsofteng/gini/inter"
	"github.com/jaredsofteng/gini/logic/aiger"
	"github.com/pkg/profile"
	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var logFile string

func main() {
	startFolder := path.Join(os.Getenv("USERPROFILE"), "Downloads")
	file1, err := dialog.File().SetStartDir(startFolder).Filter("DIMACS File", "bz2", "cnf", "gz", "aag","aig").Load()
	if err != nil {
		log.Fatal("File Selection Incomplete")
		return
	}
	ok := dialog.Message("%s", "Import entire Directory?").Title("Import Scope").YesNo()
	applyPre := dialog.Message("%s", "Apply Preprocessing?").Title("Preprocessing").YesNo()
	if ok {

		defer profile.Start().Stop()
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
			if applyPre {
				err := err
				f, err = filePreprocess(f)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
			logToFile(path.Base(f))
			g := readFile(f)
			solveMainRoutine(g)
		}
	} else {
		setLogFile(file1)
		if applyPre {
			defer profile.Start().Stop()
			startTime := time.Now()
			err := err
			file1, err = filePreprocess(file1)
			if err != nil {log.Fatal(err.Error())
				return
			}
			fileProcessTime := time.Since(startTime)
			logToFile("File Preprocessed Time (incl read) = " + fileProcessTime.String())
		}
		logToFile(path.Base(file1))
		g := readFile(file1)
		preprocess(g) // TODO: do some performance wrapping for the new pre-processing code
		file,_ := os.Create(file1 + "-pp.cnf") // Temporary to compare the CNF output
		_ = g.Write(file) // TEMPORARY
		solveMainRoutine(g)
	}
}

func solveMainRoutine(g *gini.Gini) {
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

func readAiger(filepath string) *aiger.T {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var r io.Reader
	r = f
	startTime := time.Now()

	s, err := aiger.ReadBinary(r)
	if err != nil {panic(err)}
	fileReadTime := time.Since(startTime)
	logToFile("DIMACS parsing time = " + fileReadTime.String())
	return s
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

func testRand3Cnf(vars int) int {
	g := rand3Cnf(vars)
	startMem := Tools.TotalMemUsageMB()
	doSolve := g.GoSolve()
	Tools.CpuUsagePercent(100 * time.Microsecond) // Tracks CPU percent for the next 100 microseconds
	startSolve := time.Now()
	result := doSolve.Try(300*time.Second)
	endSolve := time.Since(startSolve)
	logToFile("Solve Time = " + endSolve.String())
	cpuPercentChange := Tools.CpuUsagePercent(0) // returns difference from last cpu check
	logToFile("CPU Usage % = " + strconv.FormatFloat(cpuPercentChange, 'f', 6, 64))
	memConsumed := Tools.TotalMemUsageMB() - startMem
	logToFile("Memory Usage Total = " + strconv.FormatUint(memConsumed, 10) + "MB")
	return result
}

func rand3Cnf(vars int) inter.S {
	s := gini.NewS()
	gen.HardRand3Cnf(s, vars)
	return s
}

func filePreprocess(filepath string) (newfilepath string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("could not open %q: %v", filepath, err)
	}
	defer f.Close()
	if ! strings.HasSuffix(filepath, ".cnf") {
		return "", fmt.Errorf("invalid file format for %q", filepath)
	}
	pb, err := Preprocessor.ParseCNF(f)
	if err != nil {
		return "", fmt.Errorf("could not parse DIMACS file %q: %v", filepath, err)
	}
	pb.Preprocess()
	// write to file
	filepathNoExt := strings.TrimSuffix(filepath, path.Ext(filepath))
	file,err := os.Create(filepathNoExt + "-pp.cnf")
	if err!= nil{
		fmt.Println(err)
		return
	}
	l,err := file.WriteString(pb.CNF())
	if err!=nil{
		fmt.Println(err)
		file.Close()
		return
	}
	fmt.Println(l,"CNF file created successfully!")
	file.Close()
	return file.Name(), nil
}

func solveAiger() {
	startFolder := path.Join(os.Getenv("USERPROFILE"), "Downloads")
	file1, _ := dialog.File().SetStartDir(startFolder).Filter("DIMACS File", "bz2", "cnf", "gz", "aag","aig").Load()
	setLogFile(file1)
	g := gini.New()
	aig := readAiger(file1)
	aig.C.ToCnf(g)
	r := solveFile(g, 30*time.Second)
	printResult(r)
	_ = open.Start(logFile)
	return
}

func preprocess(g *gini.Gini) {
	pregini.Subsumption(g)
	//pregini.SelfSubsumption(g) // Performs selfsub on a gini solver
	return
}

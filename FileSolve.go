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
	file1, err := dialog.File().SetStartDir(startFolder).Filter("DIMACS File", "bz2", "cnf", "gz", "aag","aig","dimacs").Load()
	if err != nil {
		log.Fatal("File Selection Incomplete")
		return
	}
	ok := dialog.Message("%s", "Import entire Directory?").Title("Import Scope").YesNo()
	applyPre := dialog.Message("%s", "Apply Preprocessing?").Title("Preprocessing").YesNo()
	if ok {
		files, err := Tools.WalkMatch(path.Dir(file1), "*.cnf")
		bzFiles, err := Tools.WalkMatch(path.Dir(file1), "*.bz2")
		gzFiles, err := Tools.WalkMatch(path.Dir(file1), "*.gz")
		files = append(files, bzFiles...)
		files = append(files, gzFiles...)
		if err != nil {
			log.Fatal(err)
		}
		setLogDir(file1)
		for _, f := range files {
			if applyPre {
				//defer profile.Start().Stop()
				//err := err
				//f, err = filePreprocess(f)
				//if err != nil {
				//	log.Fatal(err.Error())
				//}
			}
			writeCSVtoLog(path.Base(f))
			writeCSVtoLog(time.Now().Format("2006-03-09"))
			writeCSVtoLog(time.Now().Format("15:04:05"))
			g := readFile(f)
			var rem int
			startTime := time.Now()
			if applyPre {
				rem = preprocess(g)
			}
			fileProcessTime := time.Since(startTime)
			writeCSVtoLog(fileProcessTime.String())
			writeCSVtoLog(strconv.FormatInt(int64(rem), 10))
			solveMainRoutine(g)
			logToFile("")
		}
		_ = open.Start(logFile)
	} else {
		setLogFile(file1)
		if applyPre {
			//defer profile.Start().Stop()
			//startTime := time.Now()
			//err := err
			//file1, err = filePreprocess(file1)
			//fileProcessTime := time.Since(startTime)
		}
		writeCSVtoLog(path.Base(file1))
		writeCSVtoLog(time.Now().Format("2006-03-09"))
		writeCSVtoLog(time.Now().Format("15:04:05"))
		g := readFile(file1)
		var rem int
		startTime := time.Now()
		if applyPre {
			rem = preprocess(g)
		}
		fileProcessTime := time.Since(startTime)
		writeCSVtoLog(fileProcessTime.String())
		writeCSVtoLog(strconv.FormatInt(int64(rem), 10))
		//file,_ := os.Create(file1 + "-pregini2.cnf") // Temporary to compare the CNF output
		//_ = g.Write(file) // TEMPORARY
		solveMainRoutine(g)
		_ = open.Start(logFile)
	}
}

func solveMainRoutine(g *gini.Gini) {
	maxSolveTime := time.Second * 10
	r := solveFile(g, maxSolveTime)
	printResult(r)
}

func readFile(filepath string) *gini.Gini {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	startTime := time.Now()
	g, err := readFileReader(f)
	if err != nil {panic(err)}
	fileReadTime := time.Since(startTime)
	writeCSVtoLog(fileReadTime.String())
	return g
}

func readAiger(r io.Reader) (*aiger.T, error) {
	t, err := aiger.ReadBinary(r)
	if err != nil {return nil, err}
	return t, err
}

func readFileReader(f *os.File) (*gini.Gini, error) {
	var r io.Reader
	var e error
	g := gini.New()
	switch path.Ext(f.Name()) {
	case ".AIG":
	case ".aig":
	case ".aag":
	case ".AAG":
		r = f
		t, e := readAiger(r)
		if e != nil {
			log.Fatal(e.Error())
		}
		t.C.ToCnf(g)
	case ".GZ":
	case ".gz":
		r, e = gzip.NewReader(f)
		if e != nil {
			log.Fatal(e.Error())
		}
		g, e = gini.NewDimacs(r)
		if e != nil {
			log.Fatal(e.Error())
		}
	case ".BZ2":
	case ".bz2":
		r = bzip2.NewReader(f)
		g, e = gini.NewDimacs(r)
		if e != nil {
			log.Fatal(e.Error())
		}
	case ".CNF":
	case ".cnf":
		r = f
		g, e = gini.NewDimacs(r)
		if e != nil {
			log.Fatal(e.Error())
		}
	default:
		log.Fatal("Invalid File format - must be .aig .aag .gz .bz2 or .cnf")
	}
	return g, e
}

func solveFile(g *gini.Gini, timeout time.Duration) int {
	startMem := Tools.TotalMemUsageMB()
	doSolve := g.GoSolve()
	Tools.CpuUsagePercent(100 * time.Microsecond) // Tracks CPU percent for the next 100 microseconds
	startSolve := time.Now()
	result := doSolve.Try(timeout)
	endSolve := time.Since(startSolve)
	writeCSVtoLog(endSolve.String())
	cpuPercentChange := Tools.CpuUsagePercent(0) // returns difference from last cpu check
	writeCSVtoLog(strconv.FormatFloat(cpuPercentChange, 'f', 6, 64))
	memConsumed := Tools.TotalMemUsageMB() - startMem
	writeCSVtoLog(strconv.FormatUint(memConsumed, 10) + "MB")
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

func writeCSVtoLog(s string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { log.Fatal(err)}
	if _, err := fmt.Fprint(f, s + ","); err != nil {
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
	var newFilenameExt = ".csv"
	newFile := filename[0:len(filename)-len(filenameExt)]
	logFile = newFile + "-result" + newFilenameExt
	if _, err := os.Stat(logFile); err == nil {
		logToFile("") // Newline
	} else if os.IsNotExist(err) {
		logToFile("File Name,Date,Time,DIMACS Time,Preprocess Time,Clauses Removed,Solve Time,CPU,MEM,Result,")
	}
}

func setLogDir(fileIn string) {
	newFilenameExt := ".csv"
	logFile = path.Dir(fileIn) + "Directory-CNF-Results" + newFilenameExt
	if _, err := os.Stat(logFile); err == nil {
	} else if os.IsNotExist(err) {
		logToFile("File Name,Date,Time,DIMACS Time,Preprocess Time,Clauses Removed,Solve Time,CPU,MEM,Result,")
	}
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
		resultStr = "ERR"
	}
	writeCSVtoLog(resultStr)
}

func testRand3Cnf(vars int) int {
	g := rand3Cnf(vars)
	startMem := Tools.TotalMemUsageMB()
	doSolve := g.GoSolve()
	Tools.CpuUsagePercent(100 * time.Microsecond) // Tracks CPU percent for the next 100 microseconds
	startSolve := time.Now()
	result := doSolve.Try(300*time.Second)
	endSolve := time.Since(startSolve)
	writeCSVtoLog("Solve Time = " + endSolve.String())
	cpuPercentChange := Tools.CpuUsagePercent(0) // returns difference from last cpu check
	writeCSVtoLog("CPU Usage % = " + strconv.FormatFloat(cpuPercentChange, 'f', 6, 64))
	memConsumed := Tools.TotalMemUsageMB() - startMem
	writeCSVtoLog("Memory Usage Total = " + strconv.FormatUint(memConsumed, 10) + "MB")
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

func preprocess(g *gini.Gini) int {
	clauseRem := pregini.Subsumption(g)
	//pregini.SelfSubsumption(g) // Performs selfsub on a gini solver
	return clauseRem
}

package GiniBench

import (
	"fmt"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/internal/xo"

	"io"
	"os"
	"path"
	"time"
)

func benchmarkFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r io.Reader
	r = f
	startTime := time.Now()
	g, _ := xo.NewSDimacs(r)
	fileReadTime := time.Since(startTime)
	stats := xo.NewStats()
	doSolve := g.GoSolve()
	result := doSolve.Try(time.Second*30)
	g.ReadStats(stats)
	exportResults(g)
}

func exportResults(g gini.Gini, newFile string) {
	var filename = newFile
	var filenameExt = path.Ext(filename)
	newFile = filename[0:len(filename)-len(filenameExt)]
	newFile = newFile + "-result" + filenameExt
	f, _ := os.Create(newFile)
	f.Write(g.)
}
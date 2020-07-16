package GiniBench

import (
	"fmt"
	"github.com/irifrance/gini"
	"io"
	"os"
	"path"
)

func benchmarkFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r io.Reader
	r = f
	g, _ := gini.NewDimacs(r)
	if g.Solve() == 1 {

	}
}

func exportResults(g gini.Gini, newFile string) {
	var filename = newFile
	var filenameExt = path.Ext(filename)
	newFile = filename[0:len(filename)-len(filenameExt)]
	newFile = newFile + "-result" + filenameExt
	f, _ := os.Create(newFile)
	f.Write(g.)
}
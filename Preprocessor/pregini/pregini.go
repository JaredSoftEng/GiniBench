package pregini

import (
	"fmt"
	"github.com/jaredsofteng/gini"
	"github.com/jaredsofteng/gini/z"
	"sort"
	"time"
)

// Takes as input a gini solver, and perform self-subsumption methods on the clause database.
// NOTE: I had to modify the gini source code in order to expose the Cdb (ClauseDB), see fork in JaredSoftEng
func SelfSubsumption (g *gini.Gini) {
	Clauses := g.ClauseDB()
	d := Clauses.CDat
	var ret error
	d.Forall(func(i int, o z.C, ms []z.Lit) {
		if ret != nil {return}
		for _, m := range ms {
			_, e := fmt.Print(m.Dimacs())
			if e != nil {
				ret = e
				return
			}
			fmt.Print(" ")
		}
		_, e := fmt.Println("0")
		if e != nil {ret=e}
	})
	return
}

// A very fast subsumption check using only the watched literals (2) of each clause. Could be unnecessary?
func WatchedBinarySubsumption (g *gini.Gini) []z.C {
	WatchedLits := g.ClauseDB().Vars.Watches
	var remClauses []z.C
	for iLit := 2; iLit < len(WatchedLits); iLit++ {
		if len(WatchedLits[iLit]) > 1 { // The literal occurs more than once
			Watches := WatchedLits[iLit]
			for j := 0; j < len(Watches); j++ {
				if Watches[j] >= (1 << 63) { // Compare with each other watched literal set
					Lit1 := Int2Lit(iLit)
					Lit2 := Watches[j].Other()
					iLit2 := Lit2Int(Lit2)
					wLit2 := WatchedLits[iLit2]
					if iLit2 > iLit && len(wLit2) > 1 {
						for k := 0; k < len(wLit2); k++ {
							if wLit2[k].Other() == Lit1 && wLit2[k].C() != Watches[j].C() {
								remClauses = append(remClauses, wLit2[k].C())
							}
						}
					}
				}
			}
		}
	}
	return remClauses
}

// A more thorough Subsumption mechanism using watched literals. For any watched literal, will compare the clauses to test for subsumption.
// BUT: This is a 2-literal watched schema, ie (3 4 0) would not subsume (1 2 3 4 0) as no watch is established for literals after the first two.
func WatchedSubsumption (g *gini.Gini) []z.C {
	WatchedLits := g.ClauseDB().Vars.Watches
	var remClauses []z.C
	for iLit := 2; iLit < len(WatchedLits); iLit++ {
		if len(WatchedLits[iLit]) > 1 { // The literal occurs more than once
			watchLen, cRef, cLits := FetchClauses(g, iLit)
			for j := 0; j < watchLen-1; j++ {
				currClause := cRef[j]
				for k := j+1; k < watchLen; k++ {
					nextClause := cRef[k]
					if len(cLits[j]) <= len(cLits[k]) {
						if Matches(cLits[j], cLits[k]) {
							remClauses = append(remClauses, nextClause)
						}
					} else {
						if Matches(cLits[k], cLits[j]) {
							remClauses = append(remClauses, currClause)
						}
					}
				}
			}
		}
	}
	return remClauses
}

func WatchedGiniLinear(g *gini.Gini) *gini.Gini {
	g2 := g.Copy()
	w := g2.ClauseDB().Vars.Watches
	watch := g2.ClauseDB().GetWatch()
	g2.ClauseDB().CDat.Forall(func(i int, o z.C, ms []z.Lit) {
		litSize := len(ms)
		if litSize == 1 {
			w[ms[0]] = append(w[ms[0]], watch.NewWatch(o, ms[0], false))
		}
		if litSize > 2 {
			for _, m := range ms[2:] {
				w[m] = append(w[m], watch.NewWatch(o, m, false))
			}
		}
	})
	return g2
}

func WatchedGiniFull(g *gini.Gini) *gini.Gini {
	g2 := g.Copy()
	w := g2.ClauseDB().Vars.Watches
	for i := 0; i < len(w); i++ {
		w[i] = w[i][:0] // reduces each to 0
	}
	watch := g2.ClauseDB().GetWatch()
	g2.ClauseDB().CDat.Forall(func(i int, o z.C, ms []z.Lit) {
		if len(ms) > 1 {
			for j, m := range ms {
				for k, n := range ms {
					if j == k  {
						continue
					}
					w[m] = append(w[m], watch.NewWatch(o, n, false))
				}
			}
		}
	})
	return g2
}

// Given the gini solver and an integer referring to a literal (2 = 1, 3 = -1 etc.) returns arrays containing the length of watched clauses, a pointer to it and its elements
func FetchClauses(g *gini.Gini, lit int) (int, []z.C, [][]z.Lit ) {
	var clausePointer []z.C
	var clauseLits []z.Lit
	var clauseLitArr [][]z.Lit
	watchLen := len(g.ClauseDB().Vars.Watches[lit])
	for _, w := range g.ClauseDB().Vars.Watches[lit] {
		clausePointer = append(clausePointer, w.C())
		clauseLitArr = append(clauseLitArr, g.ClauseDB().CDat.Load(w.C(), clauseLits))
		clauseLits = nil
	}
	return watchLen, clausePointer, clauseLitArr
}



// Wrapper for the Subsumption method, will compact the cDat prior to each call.
func Subsumption(g *gini.Gini) int {
	//start := time.Now()
	//cRemB := len(WatchedBinarySubsumption(g))
	//end := time.Since(start)
	//start1 := time.Now()
	//cList := WatchedSubsumption(g)
	//cRem2, _ := RemoveClauses(g, cList)
	//end1 := time.Since(start1)
	////fmt.Println("Watched Binary: " + end.String())
	//fmt.Println("Binary Watchlist Time: " + end1.String())
	//fmt.Println(cRemB)
	//fmt.Println(cRem )
	start2 := time.Now()
	g2 := WatchedGiniLinear(g)
	cList := WatchedSubsumption(g2)
	cRem2, _ := RemoveClauses(g, cList)
	end2 := time.Since(start2)
	fmt.Println("Linear Watchlist Time: " + end2.String())
	fmt.Println(cRem2)
	return cRem2
}

func deprecatedTestFunc (g *gini.Gini) {
	Clauses := g.ClauseDB()
	var rms []z.C
	rms = append(rms,1, 6, 11, 16)

//	Clauses.Remove(rms[3], rms[1], rms[0])
//	gc := Clauses.GetGC()
//	numRem, nVar := gc.CompactCDat(Clauses)
	d := Clauses.CDat
	w := Clauses.Vars.Watches[11]
 	w0 := w[0].IsBinary()

	fmt.Print(w)
	fmt.Print(w0)
	fmt.Print("Clause Len: ")
	fmt.Print(d.ClsLen)
	fmt.Print("Lit Len: ")
	fmt.Println(d.Len)
	c1 := d.Chd(2)
	fmt.Print("Clause1 Size: ")
	fmt.Println(c1.Size())
	fmt.Println(c1.String())


	var ret error
	d.Forall(func(i int, o z.C, ms []z.Lit) {
		if ret != nil {return}
		fmt.Print(o.String())
		fmt.Print("Lits: ")
		fmt.Print(d.Chd(o).Size())
		fmt.Print(" ")
		for _, m := range ms {

			_, e := fmt.Print(m.Dimacs())
			if e != nil {
				ret = e
				return
			}
			fmt.Print(" ")
		}
		fmt.Println("0")
	})
}

// First traverses clause set to unlink the associated watched literals from each clause, and adds the clause to the removal queue.
// CompactCDat is called which remaps the byte space associated with the clause set.
func RemoveClauses(g *gini.Gini, c []z.C) (int, int) {
	cLocSlice(c).Sort()
	c = uniq(c)
	var nBytesRem int
	if g.ClauseDB().CDat.CompactReady(len(c), len(c)*4) {
		g.ClauseDB().Remove(c...)
	} else {
		g.ClauseDB().Remove(c...)
		_, nBytesRem = g.ClauseDB().GetGC().CompactCDat(g.ClauseDB())
	}
	return len(c), nBytesRem
}

// Compares arrays of z.Lit (A -> B), returns true if B contains all the elements in A (for performance, assumes literal ordering )
func Matches(a []z.Lit, b []z.Lit) bool {
	for i := range a {
		if !Has(a[i], b) { // aLit not found in array
			return false
		}
	}
	return true
}

// linear search for a literal in a set of literals giving the index found
func Has(a z.Lit, l []z.Lit) bool {
	for i := range l {
		if l[i] == a {return true}
	}
	return false
}

// Returns the value of the variable index based on a z.Lit
func Lit2Int(z z.Lit) int {
	if z.Sign() == 1 {
		return z.Dimacs()*2
	} else {
		return z.Dimacs()*-2+1
	}
}

func Int2Lit(i int) z.Lit {
	if i % 2 == 0 {
		return z.Dimacs2Lit(i/2)
	} else {
		return z.Dimacs2Lit(-(i-1)/2)
	}
}

type cLocSlice []z.C

func (p cLocSlice) Len() int {
	return len(p)
}

func (p cLocSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p cLocSlice) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p cLocSlice) Sort() {
	sort.Sort(p)
}

func uniq(cs []z.C) []z.C {
	if len(cs) <= 1 {
		return cs
	}
	i := 0
	j := 1
	last := cs[0]
	var cur z.C
	N := len(cs)
	for j < N {
		cur = cs[j]
		if cur != last {
			i++
			cs[i] = cur
			last = cur
		}
		j++
	}
	return cs[:i+1]
}
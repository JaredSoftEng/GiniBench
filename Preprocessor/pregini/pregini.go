package pregini

import (
	"github.com/jaredsofteng/gini"
	"github.com/jaredsofteng/gini/z"
)

// Takes as input a gini solver, and performs self-subsumption on the clause database using binary pivot clauses.
func BinarySelfSubsumption (g *gini.Gini) (int, []z.C) {
	Clauses := g.ClauseDB()
	d := Clauses.CDat
	var remClauseMap = make(map[z.C]z.C)
	var remClauses []z.C
	remLit := 0
	d.Forall(func(i int, o z.C, ms []z.Lit) {
		if _, ok := remClauseMap[o]; !ok {
			var remLitMap = make(map[z.Lit]z.Lit)
			isClauseRem := false
			// Go through each literal in the clause
			for _, m := range ms {
				// Go through each watch in the negation of the literal
				mWatches := Clauses.Vars.Watches[m.Not()]
				for j := range mWatches {
					// If the watch is binary
					if mWatches[j].IsBinary() {
						// Check if the other literal of the binary watch is in the current clause
						if Has(mWatches[j].Other(), ms) {
							// if so, eliminate the literal
							remLitMap[m] = m
							remClauseMap[o] = o
							isClauseRem = true
						}
					}
				}
			}
			// If any literals were removed, make a new clause of the pivoted literals (in literal order)
			if isClauseRem {
				remLit += len(ms)
				for i := range ms {
					if _, ok := remLitMap[ms[i]]; !ok {
						g.Add(ms[i])
						remLit--
					}
				}
				g.Add(0)
			}
		}
	})
	for j := range remClauseMap {
		if _, ok := remClauseMap[j]; ok {
			remClauses = append(remClauses, j) // Now is ordered by z.C
		}
	}
	return remLit, remClauses
}

// A fast subsumption check using only the watched literals (2) of each clause. Does not match many clauses in practice.
func WatchedBinarySubsumption (g *gini.Gini) []z.C {
	WatchedLits := g.ClauseDB().Vars.Watches
	var remClauseMap = make(map[z.C]z.C)
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
								remClauseMap[wLit2[k].C()] = wLit2[k].C()
							}
						}
					}
				}
			}
		}
	}
	for j := range remClauseMap {
		if _, ok := remClauseMap[j]; ok {
			remClauses = append(remClauses, j) // Now is ordered by z.C
		}
	}
	return remClauses
}

// A thorough Subsumption mechanism using the watched literal list. For any watched literal, will compare the clauses to test for subsumption.
// If using the unmodified watchlist; it is a 2-literal watched schema, ie (3 4 0) would not subsume (1 2 3 4 0) as no watch is established for literals after the first two.
func WatchedSubsumption (g *gini.Gini) []z.C {
	WatchedLits := g.ClauseDB().Vars.Watches
	var remClauseMap = make(map[z.C]z.C)
	var remClauses []z.C
	for iLit := 2; iLit < len(WatchedLits); iLit++ {
		watchLen := len(WatchedLits[iLit])
		if watchLen > 1 { // The literal occurs more than once
			cRef, cLits := FetchClauses(g, iLit)
			for j := 0; j < watchLen-1; j++ {
				currClause := cRef[j]
				for k := j+1; k < watchLen; k++ {
					nextClause := cRef[k]
					if len(cLits[j]) <= len(cLits[k])  {
						if Matches(cLits[j], cLits[k]) {
							remClauseMap[nextClause] = nextClause
						}
					} else {
						if Matches(cLits[k], cLits[j]) {
							remClauseMap[currClause] = currClause
						}
					}
				}
			}
		}
	}
	for j := range remClauseMap {
		if _, ok := remClauseMap[j]; ok {
			remClauses = append(remClauses, j) // Now is ordered by z.C
		}
	}
	return remClauses
}

// A sample subsumption mechanism using the hash of clauses as a comparator to test for subsumption.
func WatchedHashSubsumption (g *gini.Gini) []z.C {
	WatchedLits := g.ClauseDB().Vars.Watches
	var remClauses []z.C
	clauseHash := CreateClauseHash(g)
	for iLit := 2; iLit < len(WatchedLits); iLit++ {
		watchLen := len(WatchedLits[iLit])
		if watchLen > 1 { // The literal occurs more than once
			cRef, cLits := FetchClauses(g, iLit)
			for j := 0; j < watchLen-1; j++ {
				currClause := cRef[j]
				for k := j+1; k < watchLen; k++ {
					nextClause := cRef[k]
					if len(cLits[j]) <= len(cLits[k]) {
						// if nextClause is the bigger clause being subsumed
						if HashCheck(clauseHash,currClause,nextClause){
							if Matches(cLits[j], cLits[k]) {
								remClauses = append(remClauses, nextClause)
							}
						}
					} else {
						// if currClause is the bigger clause being subsumed
						if HashCheck(clauseHash,nextClause,currClause) {
							if Matches(cLits[k], cLits[j]) {
								remClauses = append(remClauses, currClause)
							}
						}
					}
				}
			}
		}
	}
	return remClauses
}

// This takes a gini solver as an input, and outputs a copy of it where the watchlist has been expanded to include all the literals of every clause.
// Note: This also includes unit clauses (a single literal).
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

// Given the gini solver instance and an integer referring to a literal (2 = 1, 3 = -1 etc.) returns arrays containing clause pointers and their elements.
func FetchClauses(g *gini.Gini, lit int) ([]z.C, [][]z.Lit ) {
	var clausePointer []z.C
	var clauseLits []z.Lit
	var clauseLitArr [][]z.Lit
	//watchLen := len(g.ClauseDB().Vars.Watches[lit])
	cData := g.ClauseDB().CDat
	for _, w := range g.ClauseDB().Vars.Watches[lit] {
		clausePointer = append(clausePointer, w.C())
		clauseLitArr = append(clauseLitArr, cData.Load(w.C(), clauseLits))
		clauseLits = nil
	}
	return clausePointer, clauseLitArr
}

// Wrapper for the Subsumption method
func FullSubsumption(g *gini.Gini) (int, int) {
	g2 := WatchedGiniLinear(g)
	cList := WatchedSubsumption(g2)
	cRem, cLit := RemoveClauses(g, cList)
	return cRem, cLit
}

// Wrapper for the Subsumption method
func Subsumption(g *gini.Gini) (int, int) {
	cList := WatchedSubsumption(g)
	cRem, cLit := RemoveClauses(g, cList)
	return cRem, cLit
}

// Wrapper for the SelfSubsumption method, will compact the cDat prior to each call.
func SelfSubsumption(g *gini.Gini) (int, int) {
	remLit, cList := BinarySelfSubsumption(g)
	RemoveClauses(g, cList)
	return len(cList), remLit
}

// First traverses clause set to unlink the associated watched literals from each clause, and adds the clause to the removal queue.
// CompactCDat is called which remaps the byte space associated with the clause set.
func RemoveClauses(g *gini.Gini, c []z.C) (int, int) {
	var nBytesRem int
	if g.ClauseDB().CDat.CompactReady(len(c), len(c)*4) {
		g.ClauseDB().Remove(c...)
	} else {
		g.ClauseDB().Remove(c...)
		_, nBytesRem = g.ClauseDB().GetGC().CompactCDat(g.ClauseDB())
	}
	return len(c), nBytesRem
}

// Compares arrays of z.Lit (A -> B), returns true if B contains all the elements in A
func Matches(a []z.Lit, b []z.Lit) bool {
	for i := range a {
		if !Has(a[i], b) { // aLit not found in array
			return false
		}
	}
	return true
}

// linear search for a literal in a set of literals (no ordering assumed)
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

// Creates a 64-bit integer representing the clause
func CreateClauseHash(g *gini.Gini) map[z.C]uint64 {
	clauseHash:= make(map[z.C]uint64)
	d := g.ClauseDB().CDat
	var ret error
	val := uint64(0)
	d.Forall(func(i int, o z.C, ms []z.Lit) {

		if ret != nil {return}
		for _, m := range ms {
			valTemp := Hash(Lit2Int(m))
			// perform bitwise OR on the values
			val = valTemp | val
		}
		// add the finished clause signature
		clauseHash[o] = val
	})
	return clauseHash

}

// return "1" shifted left the number of times equal to the hash integer and casted as a 64 bit integer
func Hash(i int) uint64 {
	i = i % 63
	return uint64(1<<i)
}

// takes the bitwise AND between c1 and complement (^c2) of c2
func HashCheck(clauseHash map[z.C]uint64, i z.C, j z.C) bool{
	c1 := clauseHash[i]
	c2 := clauseHash[j]
	if c1 & (^c2) != uint64(0){
		return false
	} else{
		return true
	}
}
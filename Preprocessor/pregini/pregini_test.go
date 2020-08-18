package pregini

import (
	"github.com/jaredsofteng/gini/z"
	"testing"
)

func TestLit2Int(t *testing.T) {
	for i := 1; i < 100; i++ {
		if Lit2Int(z.Dimacs2Lit(i)) != 2*i {
			t.Errorf("Positive Lit mismatch")
		}
		if Lit2Int(z.Dimacs2Lit(-i)) != 2*i+1 {
			t.Errorf("Negative Lit mismatch")
		}
	}
}

func TestHas(t *testing.T) {
	var litArr []z.Lit
	for i := 1; i < 10; i++ {
		litArr = append(litArr, z.Dimacs2Lit(i))
	}
	litCompare := z.Dimacs2Lit(4)
	if Has(litCompare, litArr) == -1 {
		t.Errorf("Has function invalid")
	}
}

func TestMatches(t *testing.T) {
	var litArr []z.Lit
	var litArr2 []z.Lit
	for i := 1; i < 10; i++ {
		litArr = append(litArr, z.Dimacs2Lit(i))
	}
	for i := 3; i < 6; i++ {
		litArr2 = append(litArr2, z.Dimacs2Lit(i))
	}
	if !Matches(litArr2, litArr) {
		t.Errorf("Matches function invalid")
	}
}
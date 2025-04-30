package rw

import (
	"strings"
	"testing"

	"github.com/vic/eff.go"
)

func TestReadWrite(t *testing.T) {
	type S []string

	e := eff.FlatMap(Read[S](), func(s *S) eff.Eff[WriteAb[S], WriteRs[S]] {
		n := append(*s, "world")
		e := Write(&n)
		return e
	})

	st := &S{"hello"}
	rh := ReadHandler(func() *S { return st })
	wh := WriteHandler(func(s *S) { st = s })
	x := eff.ProvideBoth(e, rh.Ability(), wh.Ability())

	eff.Eval(x)

	s := strings.Join(*st, " ")
	if s != "hello world" {
		t.Errorf("unexpected value %v", *st)
	}
}

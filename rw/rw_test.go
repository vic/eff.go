package rw

import (
	"strings"
	"testing"

	"github.com/vic/eff.go"
)

func TestReadWrite(t *testing.T) {
	type S []string

	e := eff.FlatMap(Read[S](), func(s *S) WriteEff[S] {
		n := append(*s, "world")
		e := Write(&n)
		return e
	})

	st := &S{"hello"}
	rh := ReadHandler(func() *S { return st })
	wh := WriteHandler(func(s *S) { st = s })
	e1 := eff.ProvideLeft(e, rh.Ability())
	e2 := eff.Provide(e1, wh.Ability())

	_, err := eff.Eval(e2)
	if err != nil {
		t.Error(err)
		return
	}

	s := strings.Join(*st, " ")
	if s != "hello world" {
		t.Errorf("unexpected value %v", *st)
	}
}

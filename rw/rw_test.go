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
	e1 := eff.HandleBoth(rh, wh, e)

	_, err := eff.Eval(e1)
	if err != nil {
		t.Error(err)
		return
	}

	s := strings.Join(*st, " ")
	if s != "hello world" {
		t.Errorf("unexpected value %v", *st)
	}
}

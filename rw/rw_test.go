package rw

import (
	"strings"
	"testing"

	fx "github.com/vic/fx.go"
)

func TestReadWrite(t *testing.T) {
	type S []string
	type E = fx.Fx[fx.And[ReadAb[S], WriteAb[S]], fx.Nil]

	var e E = fx.FlatMap(Read[S](), func(s *S) WriteFx[S, fx.Nil] {
		n := append(*s, "world")
		return Write(&n)
	})

	st := &S{"hello"}
	rh := ReadHandler(func() *S { return st })
	wh := WriteHandler(func(s *S) { st = s })
	x := fx.AndCollapse(fx.ProvideAB(e, rh, wh))
	fx.Eval(x)

	s := strings.Join(*st, " ")
	if s != "hello world" {
		t.Errorf("unexpected value %v", *st)
	}
}

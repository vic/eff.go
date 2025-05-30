package rw

import (
	"strings"
	"testing"

	"github.com/vic/fx.go/fx"
)

func TestReadWrite(t *testing.T) {
	type S []string
	type E = fx.Fx[fx.And[ReadAb[S], WriteAb[S]], fx.Nil]

	var e E = fx.FlatMap(Read[S](), func(s *S) WriteFx[S, fx.Nil] {
		n := append(*s, "world")
		return Write(&n)
	})

	st := &S{"hello"}
	rh := ReadService(func() *S { return st })
	wh := WriteService(func(s *S) { st = s })
	x := fx.AndCollapse(fx.ProvideFirsts(e, rh, wh))
	fx.Eval(x)

	s := strings.Join(*st, " ")
	if s != "hello world" {
		t.Errorf("unexpected value %v", *st)
	}
}

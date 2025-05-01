package abort

import (
	"testing"

	"github.com/vic/eff.go"
)

func TestSuccess(t *testing.T) {
	type Err = string
	value := 22
	e := eff.Value[AbortAb[Err, int]](&value)
	x := HandleAbort(e)
	var r *Result[Err, int] = eff.Eval(x)
	val, err := (*r)()
	if err != nil {
		t.Error(*err)
	}
	if *val != value {
		t.Errorf("unexpected result %v", *val)
	}
}

func TestFailure(t *testing.T) {
	t.SkipNow()
	type Err = string
	e := eff.Map(Abort[Err, int]("ahhhh"), func(_ Result[Err, int]) int {
		panic("unreachable")
	})
	x := HandleAbort(e)
	var r *Result[Err, int] = eff.Eval(x)
	val, err := (*r)()
	if *err != "ahhhh" {
		t.Errorf("unexpected err %v", *err)
	}
	if val != nil {
		t.Errorf("unexpected success %v", *val)
	}
}

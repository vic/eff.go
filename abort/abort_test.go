package abort

import (
	"testing"

	"github.com/vic/fx.go/fx"
)

func TestSuccess(t *testing.T) {
	type Ok = int
	type Err = string
	value := 22
	e := Succeed[Err](22)
	// Applying the handler directly
	x := AbortHandler(e)
	var r Result[Ok, Err] = fx.Eval(x)
	val, err := r()
	if err != nil {
		t.Error(*err)
	}
	if *val != value {
		t.Errorf("unexpected result %v", *val)
	}
}

func TestFailure(t *testing.T) {
	type Ok = int
	type Err = string
	e := fx.Map(Abort[Ok]("ahhhh"), func(_ Ok) int {
		panic("BUG: mapping on aborted eff should be unreachable")
	})
	// Another way of applying the abort handler.
	x := fx.ProvideLeft(fx.Suspend[ResultHn[Ok, Err]](e), AbortHandler[Ok, Err])
	var r Result[Ok, Err] = fx.Eval(x)
	val, err := r()
	if *err != "ahhhh" {
		t.Errorf("unexpected err %v", *err)
	}
	if val != nil {
		t.Errorf("unexpected success %v", *val)
	}
}

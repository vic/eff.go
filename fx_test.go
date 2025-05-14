package fx

import (
	"testing"
)

type Person struct {
	name string
}

func TestPure(t *testing.T) {
	person := Person{}
	e := Pure(&person)
	found := Eval(e)
	person.name = "ulrika"
	if found.name != "ulrika" {
		t.Errorf("not the person I want: %v", person)
	}
}

func TestFunc(t *testing.T) {
	strLen := func(s string) int {
		return len(s)
	}
	e := Func(strLen)
	provided := Provide(e, "hello")
	result := Eval(provided)
	if result != len("hello") {
		t.Errorf("unexpected result %v", result)
	}
}

type printRq = func(string) FxPure[int]
type printHn = func(Fx[And[printRq, Nil], int]) Fx[Nil, int]
type printAb = And[printHn, Nil]
type printFx = Fx[printAb, int]

func PrintLn(line string) printFx {
	return Request[printHn](line)
}

func dontPrintHandler() printHn {
	return Handler(func(line string) FxPure[int] {
		r := len(line)
		return Pure(r)
	})
}

func TestHandleSimple(t *testing.T) {
	e := PrintLn("hello")
	h := dontPrintHandler()
	f := ProvideLeft(e, h)
	v := Eval(f)
	if v != len("hello") {
		t.Errorf("unexpected value %v", v)
	}
}

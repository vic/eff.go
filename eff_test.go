package eff

import (
	"testing"
)

type Person struct {
	name string
}

func TestPure(t *testing.T) {
	person := Person{}
	e := Pure(&person)
	found, err := Eval(e)
	person.name = "ulrika"
	if err != nil {
		t.Error(err)
		return
	}
	if found.name != "ulrika" {
		t.Errorf("not the person I want: %v", person)
	}
}

type Console[S any] struct{}
type printRq string
type printRs int
type printAb[S any] = Ability[printRq, printRs, S]
type printEff[S any] = Eff[printAb[S], printRs]

func (c Console[S]) PrintLn(line string) printEff[S] {
	return Suspend[printEff[S]](printRq(line))
}

func dontPrintHandler() Handler[printRq, printRs, Nil] {
	return func(q printRq, f Cont[Nil, printRs]) Eff[Nil, printRs] {
		r := len(q)
		return f(printRs(r))
	}
}

func TestHandleSimple(t *testing.T) {
	h := dontPrintHandler()
	e := Console[Nil]{}.PrintLn("hello")
	f := Provide(e, h.Ability())
	v, err := Eval(f)
	if err != nil {
		t.Error(err)
		return
	}
	if int(*v) != len("hello") {
		t.Errorf("unexpected value %v", *v)
	}
}

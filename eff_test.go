package eff

import "testing"

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

type Console[printSt any] struct{}
type printRq string
type printRs int
type printEff[printSt any] = SusEff[printRq, printRs, printSt]

func (c Console[printSt]) PrintLn(line string) printEff[printSt] {
	return Suspend[printEff[printSt]](printRq(line))
}

func dontPrintHandler() Handler[printRq, printRs, None] {
	return func(q printRq, f Cont[None, printRs]) Eff[None, printRs] {
		r := len(q)
		return f(printRs(r))
	}
}

func TestHandleSimple(t *testing.T) {
	h := dontPrintHandler()
	e := Console[None]{}.PrintLn("hello")
	f := Handle(e, h)
	v, err := Eval(f)
	if err != nil {
		t.Error(err)
		return
	}
	if int(*v) != len("hello") {
		t.Errorf("unexpected value %v", *v)
	}
}

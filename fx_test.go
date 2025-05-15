package fx

import (
	"strconv"
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

func TestApply(t *testing.T) {
	type F = func(string) int
	var strLen F = func(s string) int {
		return len(s)
	}
	var e Fx[F, int] = Apply[F]("hello")
	provided := Provide(e, strLen)
	result := Eval(provided)
	if result != len("hello") {
		t.Errorf("invalid result %v", result)
	}
}

type printFn = func(string) FxPure[int]

func PrintLn(line string) Fx[And[printFn, Nil], int] {
	return Suspend[printFn](line)
}

func printLn(line string) FxPure[int] {
	return Pure(len(line))
}

func TestHandleSimple(t *testing.T) {
	e := Map(PrintLn("hello"), strconv.Itoa)
	f := ProvideLeft(e, printLn)
	var v string = Eval(f)
	if v != strconv.Itoa(len("hello")) {
		t.Errorf("unexpected value %v", v)
	}
}

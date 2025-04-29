package eff

import "fmt"

type immediate[V any] *V
type suspended[S, V any] func(S) Eff[S, V]
type Eff[S, V any] func() (immediate[V], suspended[S, V])

type None *None

func Pure[V any](v *V) Eff[None, V] {
	return Value[None](v)
}

func Value[S, V any](v *V) Eff[S, V] {
	return func() (immediate[V], suspended[S, V]) {
		return v, nil
	}
}

func Func[S, V any](f func(S) V) Eff[S, V] {
	return func() (immediate[V], suspended[S, V]) {
		return nil, func(s S) Eff[S, V] {
			v := f(s)
			return Value[S](&v)
		}
	}
}

func Id[V any](v V) V { return v }

func Ctx[V any]() Eff[V, V] {
	return Func(Id[V])
}

type And[A, B any] func() (*A, *B)

func left[B, A any](a *A) And[A, B] {
	return func() (*A, *B) { return a, nil }
}

func right[A, B any](b *B) And[A, B] {
	return func() (*A, *B) { return nil, b }
}

func both[N And[A, B], A, B any](n N) (*A, *B) {
	return n()
}

func fst[A, B any](p And[A, B]) A {
	a, _ := p()
	return *a
}

func AndNone[S, V any](e Eff[S, V]) Eff[And[S, None], V] {
	return cont(fst, func(v immediate[V]) Eff[And[S, None], V] {
		return Value[And[S, None]](v)
	})(e)
}

func Provide[S, V any](e Eff[S, V], s S) Eff[None, V] {
	return ProvideLeft(AndNone(e), s)
}

func ProvideBoth[A, B, V any](e Eff[And[A, B], V], a A, b B) Eff[None, V] {
	return Provide(ProvideLeft(e, a), b)
}

func ProvideRight[A, B, V any](e Eff[And[A, B], V], b B) Eff[A, V] {
	return ProvideLeft(Rotate(e), b)
}

func ProvideLeft[A, B, V any](e Eff[And[A, B], V], a A) Eff[B, V] {
	imm, sus := e()
	if imm != nil {
		return Value[B](imm)
	}
	e = sus(left[B](&a))
	imm, sus = e()
	if imm != nil {
		return Value[B](imm)
	}
	return func() (immediate[V], suspended[B, V]) {
		return nil, func(b B) Eff[B, V] {
			e = sus(right[A](&b))
			return ProvideLeft(e, a)
		}
	}
}

func Rotate[A, B, V any](e Eff[And[A, B], V]) Eff[And[B, A], V] {
	rot := func(n And[B, A]) And[A, B] {
		b, a := n()
		return func() (*A, *B) { return a, b }
	}
	return cont(rot, func(v immediate[V]) Eff[And[B, A], V] {
		return Value[And[B, A]](v)
	})(e)
}

func cont[T, U, S, V any](f func(T) S, g func(immediate[V]) Eff[T, U]) func(Eff[S, V]) Eff[T, U] {
	return func(e Eff[S, V]) Eff[T, U] {
		return func() (immediate[U], suspended[T, U]) {
			immV, susV := e()
			if immV != nil {
				eff := g(immV)
				return eff()
			}
			return nil, func(t T) Eff[T, U] {
				s := f(t)
				v := susV(s)
				return cont(f, g)(v)
			}
		}
	}
}

func ContraMap[I, O, V any](e Eff[O, V], f func(I) O) Eff[I, V] {
	return cont(f, func(v immediate[V]) Eff[I, V] {
		return Value[I](v)
	})(e)
}

func Map[U, S, V any](e Eff[S, V], f func(V) U) Eff[S, U] {
	return cont(Id, func(v immediate[V]) Eff[S, U] {
		u := f(*v)
		return Value[S](&u)
	})(e)
}

func MapM[A, U, V any](e Eff[A, U], f func(U) Eff[A, V]) Eff[A, V] {
	return cont(Id[A], func(u immediate[U]) Eff[A, V] {
		v := f(*u)
		return cont(Id[A], func(v immediate[V]) Eff[A, V] {
			return Value[A](v)
		})(v)
	})(e)
}

func FlatMap[A, U, B, V any](e Eff[A, U], f func(U) Eff[B, V]) Eff[And[A, B], V] {
	return flatMap(e, f, both)
}

func flatMap[N, A, U, B, V any](e Eff[A, U], f func(U) Eff[B, V], pair func(N) (*A, *B)) Eff[N, V] {
	fst := func(n N) A {
		a, _ := pair(n)
		return *a
	}
	snd := func(n N) B {
		_, b := pair(n)
		return *b
	}
	return cont(fst, func(u immediate[U]) Eff[N, V] {
		v := f(*u)
		return cont(snd, func(v immediate[V]) Eff[N, V] {
			return Value[N](v)
		})(v)
	})(e)
}

type Cont[S, O any] func(O) Eff[S, O]
type Handler[I, O, S any] func(I, Cont[S, O]) Eff[S, O]
type Ability[I, O, S any] And[Handler[I, O, S], S]

func Suspend[E Eff[A, O], A Ability[I, O, S], I, O, S any](input I) E {
	type H = Handler[I, O, S]
	continuation := func(o O) Eff[S, O] {
		return Value[S](&o)
	}
	eff := flatMap(Ctx[H](), func(h H) Eff[S, O] {
		return h(input, continuation)
	}, func(n A) (*Handler[I, O, S], *S) {
		return n()
	})
	return E(eff)
}

func (h *Handler[I, O, S]) Ability() Ability[I, O, S] {
	return Ability[I, O, S](left[S](h))
}

func Eval[V any](e Eff[None, V]) (*V, error) {
	imm, _ := e()
	if imm != nil {
		return imm, nil
	}
	return nil, fmt.Errorf("unhandled eff of type %T", e)
}

package fx

type immediate[V any] func() V
type suspended[S, V any] func(S) Fx[S, V]
type Fx[S, V any] func() (immediate[V], suspended[S, V])

type FxPure[V any] = Fx[Nil, V]
type FxNil = FxPure[Nil]

type Nil pnil
type pnil struct{}

var PNil Nil = Nil(pnil{})
var PureNil FxNil = Pure(PNil)

func Pure[V any](v V) FxPure[V] {
	return value[Nil](v)
}

func value[S, V any](v V) Fx[S, V] {
	return func() (immediate[V], suspended[S, V]) {
		return func() V { return v }, nil
	}
}

func Func[S, V any](f func(S) V) Fx[S, V] {
	return func() (immediate[V], suspended[S, V]) {
		return nil, func(s S) Fx[S, V] {
			return func() (immediate[V], suspended[S, V]) {
				v := func() V { return f(s) }
				return v, nil
			}
		}
	}
}

func Ctx[V any]() Fx[V, V] {
	return Func(func(v V) V { return v })
}

// An effect that will never be continued.
func Halt[S, V any]() Fx[S, V] {
	return nil
}

// Replace with y if x is already Halted. Otherwise x continues.
func Replace[S, V any](y func() Fx[S, V]) func(Fx[S, V]) Fx[S, V] {
	return func(x Fx[S, V]) Fx[S, V] {
		if x == nil {
			return y()
		}
		imm, sus := x()
		if imm != nil {
			return value[S](imm())
		}
		return func() (immediate[V], suspended[S, V]) {
			return nil, func(s S) Fx[S, V] {
				return Replace(y)(sus(s))
			}
		}
	}
}

func cont[T, U, S, V any](r func(T) S, f func(immediate[V]) Fx[T, U]) func(Fx[S, V]) Fx[T, U] {
	return func(e Fx[S, V]) Fx[T, U] {
		if e == nil {
			return nil
		}
		imm, sus := e()
		if imm != nil {
			return f(imm)
		}
		return func() (immediate[U], suspended[T, U]) {
			return nil, func(t T) Fx[T, U] {
				return cont(r, f)(sus(r(t)))
			}
		}
	}
}

type And[A, B any] func() (A, B)

func ContraMap[V, S, R any](f func(R) S) func(Fx[S, V]) Fx[R, V] {
	return func(e Fx[S, V]) Fx[R, V] {
		if e == nil {
			return nil
		}
		imm, sus := e()
		if imm != nil {
			return value[R](imm())
		}
		return func() (immediate[V], suspended[R, V]) {
			return nil, func(r R) Fx[R, V] {
				s := f(r)
				return ContraMap[V](f)(sus(s))
			}
		}
	}
}

func MapM[S, U, V any](f func(U) Fx[S, V]) func(Fx[S, U]) Fx[S, V] {
	id := func(s S) S { return s }
	return cont(id, func(u immediate[U]) Fx[S, V] {
		v := f(u())
		return cont(id, func(v immediate[V]) Fx[S, V] {
			return value[S](v())
		})(v)
	})
}

func MapT[S, V, U any](f func(V) U) func(Fx[S, V]) Fx[S, U] {
	return MapM(func(v V) Fx[S, U] {
		u := f(v)
		return value[S](u)
	})
}

func Map[S, V, U any](e Fx[S, V], f func(V) U) Fx[S, U] {
	return MapT[S](f)(e)
}

func FlatMap[A, U, B, V any](e Fx[A, U], f func(U) Fx[B, V]) Fx[And[A, B], V] {
	return FlatMapT[A](f)(e)
}

func FlatMapT[A, U, B, V any](f func(U) Fx[B, V]) func(Fx[A, U]) Fx[And[A, B], V] {
	a := func(ab And[A, B]) A {
		a, _ := ab()
		return a
	}
	b := func(ab And[A, B]) B {
		_, b := ab()
		return b
	}
	return cont(a, func(u immediate[U]) Fx[And[A, B], V] {
		return cont(b, func(v immediate[V]) Fx[And[A, B], V] {
			return value[And[A, B]](v())
		})(f(u()))
	})
}

func AndNil[S, V any](e Fx[S, V]) Fx[And[S, Nil], V] {
	fst := func(n And[S, Nil]) S {
		s, _ := n()
		return s
	}
	return cont(fst, func(v immediate[V]) Fx[And[S, Nil], V] {
		return value[And[S, Nil]](v())
	})(e)
}

func AndSwap[A, B, V any](e Fx[And[A, B], V]) Fx[And[B, A], V] {
	swp := func(ba And[B, A]) And[A, B] {
		var ab And[A, B] = func() (A, B) {
			b, a := ba()
			return a, b
		}
		return ab
	}
	return cont(swp, func(v immediate[V]) Fx[And[B, A], V] {
		return value[And[B, A]](v())
	})(e)
}

func AndJoin[A, B, V any](e Fx[A, Fx[B, V]]) Fx[And[A, B], V] {
	return FlatMap(e, func(f Fx[B, V]) Fx[B, V] { return f })
}

func AndDisjoin[A, B, V any](e Fx[And[A, B], V]) Fx[A, Fx[B, V]] {
	return func() (immediate[Fx[B, V]], suspended[A, Fx[B, V]]) {
		return nil, func(a A) Fx[A, Fx[B, V]] {
			x := ProvideLeft(e, a)
			return value[A](x)
		}
	}
}

func AndCollapse[A, V any](e Fx[And[A, A], V]) Fx[A, V] {
	return func() (immediate[V], suspended[A, V]) {
		return nil, func(a A) Fx[A, V] {
			return ProvideLeft(e, a)
		}
	}
}

func ProvideA[A, B, C, V any](e Fx[And[And[A, C], B], V], a A) Fx[And[C, B], V] {
	return AndJoin(ProvideLeft(AndDisjoin(e), a))
}

func ProvideB[A, B, D, V any](e Fx[And[A, And[B, D]], V], b B) Fx[And[A, D], V] {
	return AndSwap(ProvideA(AndSwap(e), b))
}

func ProvideAB[A, B, C, D, V any](e Fx[And[And[A, C], And[B, D]], V], a A, b B) Fx[And[C, D], V] {
	return ProvideB(ProvideA(e, a), b)
}

func Provide[S, V any](e Fx[S, V], s S) Fx[Nil, V] {
	return ProvideLeft(AndNil(e), s)
}

func ProvideBoth[A, B, V any](e Fx[And[A, B], V], a A, b B) Fx[Nil, V] {
	return Provide(ProvideLeft(e, a), b)
}

func ProvideRight[A, B, V any](e Fx[And[A, B], V], b B) Fx[A, V] {
	return ProvideLeft(AndSwap(e), b)
}

func ProvideLeft[A, B, V any](e Fx[And[A, B], V], a A) Fx[B, V] {
	if e == nil {
		return nil
	}
	imm, sus := e()
	if imm != nil {
		return value[B](imm())
	}
	return func() (immediate[V], suspended[B, V]) {
		return nil, func(b B) Fx[B, V] {
			var ab And[A, B] = func() (A, B) { return a, b }
			for {
				e = sus(ab)
				if e == nil {
					return nil
				}
				imm, sus = e()
				if imm != nil {
					return value[B](imm())
				}
			}
		}
	}
}

func Apply[F ~func(I) O, I, O any](i I) Fx[F, O] {
	return Map(Ctx[F](), func(f F) O { return f(i) })
}

func Const[S, V any](v V) Fx[S, V] {
	return Map(Ctx[S](), func(_ S) V { return v })
}

func Handle[A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[A, B], O] {
	return AndJoin(Apply[A](i))
}

func Handler[A ~func(I) Fx[B, O], B, I, O any](a A) func(Fx[And[A, B], O]) Fx[B, O] {
	return func(e Fx[And[A, B], O]) Fx[B, O] { return ProvideLeft(e, a) }
}

func Suspend[F ~func(Fx[And[A, B], O]) Fx[B, O], A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[F, B], O] {
	return Handle[F](Handle[A](i))
}

func Eval[V any](e Fx[Nil, V]) V {
	for {
		if e == nil {
			panic("tried to evaluate halted effect. try using fx.Replace with another effect.")
		}
		imm, sus := e()
		if imm != nil {
			v := imm()
			return v
		}
		e = sus(PNil)
	}
}

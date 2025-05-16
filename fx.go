package fx

type Fx[S, V any] struct {
	imm func() V
	sus func(S) Fx[S, V]
	hlt func()
}

type FxPure[V any] = Fx[Nil, V]
type FxNil = FxPure[Nil]

type Nil pnil
type pnil struct{}

var PNil Nil = Nil(pnil{})
var PureNil FxNil = Pure(PNil)

func Pure[V any](v V) FxPure[V] { return Const[Nil](v) }

func identity[V any](v V) V { return v }

func Const[S, V any](v V) Fx[S, V] { return Fx[S, V]{imm: func() V { return v }} }

func Pending[S, V any](f func(S) Fx[S, V]) Fx[S, V] { return Fx[S, V]{sus: f} }

func Func[S, V any](f func(S) V) Fx[S, V] {
	return Pending(func(s S) Fx[S, V] { return Const[S](f(s)) })
}

func Ctx[V any]() Fx[V, V] { return Func(identity[V]) }

// An effect that will never be continued.
func Halt[S, V any]() Fx[S, V] { return Fx[S, V]{hlt: func() {}} }

// Replace with y if x is already Halted. Otherwise x continues.
func Replace[S, V any](y func() Fx[S, V]) func(Fx[S, V]) Fx[S, V] {
	return func(x Fx[S, V]) Fx[S, V] {
		if x.hlt != nil {
			return y()
		}
		if x.imm != nil {
			return Const[S](x.imm())
		}
		return Pending(func(s S) Fx[S, V] { return Replace(y)(x.sus(s)) })
	}
}

// Continue an effect by transforming its immediate value into another effect.
func Cont[T, U, S, V any](cmap func(T) S, fmap func(V) Fx[T, U]) func(Fx[S, V]) Fx[T, U] {
	return func(e Fx[S, V]) Fx[T, U] {
		if e.hlt != nil {
			return Halt[T, U]()
		}
		if e.imm != nil {
			return fmap(e.imm())
		}
		return Pending(func(t T) Fx[T, U] { return Cont(cmap, fmap)(e.sus(cmap(t))) })
	}
}

type And[A, B any] func() (A, B)

func swap[A, B any](ab And[A, B]) And[B, A] {
	var ba And[B, A] = func() (B, A) {
		a, b := ab()
		return b, a
	}
	return ba
}

func left[A, B any](ab And[A, B]) A {
	a, _ := ab()
	return a
}

func right[A, B any](ab And[A, B]) B {
	_, b := ab()
	return b
}

func ContraMap[V, S, R any](f func(R) S) func(Fx[S, V]) Fx[R, V] {
	return func(e Fx[S, V]) Fx[R, V] {
		if e.hlt != nil {
			return Halt[R, V]()
		}
		if e.imm != nil {
			return Const[R](e.imm())
		}
		return Pending(func(r R) Fx[R, V] {
			return ContraMap[V](f)(e.sus(f(r)))
		})
	}
}

func MapM[S, U, V any](f func(U) Fx[S, V]) func(Fx[S, U]) Fx[S, V] {
	return FlatCont(identity, identity, f)
}

func MapH[S, V, U any](f func(V) U) func(Fx[S, V]) Fx[S, U] {
	return MapM(func(v V) Fx[S, U] { return Const[S](f(v)) })
}

func Map[S, V, U any](e Fx[S, V], f func(V) U) Fx[S, U] {
	return MapH[S](f)(e)
}

func FlatMap[A, U, B, V any](e Fx[A, U], f func(U) Fx[B, V]) Fx[And[A, B], V] {
	return FlatMapH[A](f)(e)
}

func FlatMapH[A, U, B, V any](f func(U) Fx[B, V]) func(Fx[A, U]) Fx[And[A, B], V] {
	return FlatCont[And[A, B]](left, right, f)
}

func FlatCont[N, A, U, B, V any](amap func(N) A, bmap func(N) B, fmap func(U) Fx[B, V]) func(Fx[A, U]) Fx[N, V] {
	return Cont(amap, func(u U) Fx[N, V] { return Cont(bmap, Const[N, V])(fmap(u)) })
}

func AndNil[S, V any](e Fx[S, V]) Fx[And[S, Nil], V] {
	return Cont[And[S, Nil], V](left, Const)(e)
}

func AndSwap[A, B, V any](e Fx[And[A, B], V]) Fx[And[B, A], V] {
	return Cont[And[B, A], V](swap[B, A], Const)(e)
}

func AndJoin[A, B, V any](e Fx[A, Fx[B, V]]) Fx[And[A, B], V] {
	return FlatMap(e, identity)
}

func AndDisjoin[A, B, V any](e Fx[And[A, B], V]) Fx[A, Fx[B, V]] {
	return Pending(func(a A) Fx[A, Fx[B, V]] { return Const[A](ProvideLeft(e, a)) })
}

func AndCollapse[A, V any](e Fx[And[A, A], V]) Fx[A, V] {
	return Pending(func(a A) Fx[A, V] { return ProvideLeft(e, a) })
}

func ProvideFirstLeft[A, B, C, V any](e Fx[And[And[A, C], B], V], a A) Fx[And[C, B], V] {
	return AndJoin(ProvideLeft(AndDisjoin(e), a))
}

func ProvideFirstRight[A, B, D, V any](e Fx[And[A, And[B, D]], V], b B) Fx[And[A, D], V] {
	return AndSwap(ProvideFirstLeft(AndSwap(e), b))
}

func ProvideFirsts[A, B, C, D, V any](e Fx[And[And[A, C], And[B, D]], V], a A, b B) Fx[And[C, D], V] {
	return ProvideFirstRight(ProvideFirstLeft(e, a), b)
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
	if e.hlt != nil {
		return Halt[B, V]()
	}
	if e.imm != nil {
		return Const[B](e.imm())
	}
	return Pending(func(b B) Fx[B, V] {
		var ab And[A, B] = func() (A, B) { return a, b }
		for {
			e = e.sus(ab)
			if e.hlt != nil {
				return Halt[B, V]()
			}
			if e.imm != nil {
				return Const[B](e.imm())
			}
		}
	})
}

func Apply[F ~func(I) O, I, O any](i I) Fx[F, O] {
	return Map(Ctx[F](), func(f F) O { return f(i) })
}

func Suspend[A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[A, B], O] {
	return AndJoin(Apply[A](i))
}

func Handler[A ~func(I) Fx[B, O], B, I, O any](a A) func(Fx[And[A, B], O]) Fx[B, O] {
	return func(e Fx[And[A, B], O]) Fx[B, O] { return ProvideLeft(e, a) }
}

func Handle[F ~func(Fx[And[A, B], O]) Fx[B, O], A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[F, B], O] {
	return Suspend[F](Suspend[A](i))
}

func Eval[V any](e Fx[Nil, V]) V {
	for {
		if e.hlt != nil {
			panic("tried to evaluate halted effect. try using fx.Replace with another effect.")
		}
		if e.imm != nil {
			return e.imm()
		}
		e = e.sus(PNil)
	}
}

package fx

func ContraMap[V, S, R any](f func(R) S) func(Fx[S, V]) Fx[R, V] {
	return func(e Fx[S, V]) Fx[R, V] {
		if e.res != nil {
			return Stop(func() Fx[R, V] { return ContraMap[V](f)(e.res()) })
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
	return Then(amap, func(u U) Fx[N, V] { return Then(bmap, Const[N, V])(fmap(u)) })
}

package fx

func ProvideFirstLeft[A, B, C, V any](e Fx[And[And[A, C], B], V], a A) Fx[And[C, B], V] {
	return AndFlat(ProvideLeft(AndNest(e), a))
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
	if e.res != nil {
		return Stop(func() Fx[B, V] { return ProvideLeft(e.res(), a) })
	}
	if e.imm != nil {
		return Const[B](e.imm())
	}
	return Pending(func(b B) Fx[B, V] {
		var ab And[A, B] = func() (A, B) { return a, b }
		var loop func(e Fx[And[A, B], V]) Fx[B, V]
		loop = func(e Fx[And[A, B], V]) Fx[B, V] {
			for {
				e = e.sus(ab)
				if e.res != nil {
					return Stop(func() Fx[B, V] { return loop(e.res()) })
				}
				if e.imm != nil {
					return Const[B](e.imm())
				}
			}
		}
		return loop(e)
	})
}

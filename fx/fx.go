package fx

type Fx[S, V any] struct {
	// an immediate value
	imm func() V
	// the continuation of a suspended effect
	sus func(S) Fx[S, V]
	// the resume function of an stopped effect
	res func() Fx[S, V]
}

func Pure[V any](v V) FxPure[V] { return Const[Nil](v) }

func identity[V any](v V) V { return v }

func Const[S, V any](v V) Fx[S, V] { return Fx[S, V]{imm: func() V { return v }} }

func Pending[S, V any](f func(S) Fx[S, V]) Fx[S, V] { return Fx[S, V]{sus: f} }

func Func[S, V any](f func(S) V) Fx[S, V] {
	return Pending(func(s S) Fx[S, V] { return Const[S](f(s)) })
}

func Ctx[V any]() Fx[V, V] { return Func(identity[V]) }

// Continue an effect by transforming its immediate value into another effect.
func Then[T, U, S, V any](cmap func(T) S, fmap func(V) Fx[T, U]) func(Fx[S, V]) Fx[T, U] {
	return func(e Fx[S, V]) Fx[T, U] {
		if e.res != nil {
			return Stop(func() Fx[T, U] { return Then(cmap, fmap)(e.res()) })
		}
		if e.imm != nil {
			return fmap(e.imm())
		}
		return Pending(func(t T) Fx[T, U] { return Then(cmap, fmap)(e.sus(cmap(t))) })
	}
}

func Eval[V any](e Fx[Nil, V]) V {
	for {
		if e.res != nil {
			panic("tried to evaluate an stopped effect. try using fx.Resume or fx.Replace on it.")
		}
		if e.imm != nil {
			return e.imm()
		}
		e = e.sus(PNil)
	}
}

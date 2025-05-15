package abort

import fx "github.com/vic/fx.go"

type AbortFn[E any] func(E) fx.FxNil
type AbortAb[E any] = fx.And[AbortFn[E], fx.Nil]
type AbortFx[E, V any] = fx.Fx[AbortAb[E], V]

func Abort[V, E any](e E) AbortFx[E, V] {
	return fx.Map(fx.Handle[AbortFn[E]](e), func(_ fx.Nil) V {
		panic("unhandled abort effect")
	})
}

func Succeed[E, V any](v V) AbortFx[E, V] {
	return fx.Const[AbortAb[E]](v)
}

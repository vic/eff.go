package abort

import fx "github.com/vic/fx.go"

type Result[V, E any] func() (*V, *E)

func success[V, E any](v V) Result[V, E] { return func() (*V, *E) { return &v, nil } }
func failure[V, E any](e E) Result[V, E] { return func() (*V, *E) { return nil, &e } }

type AbortFn[V, E any] = func(E) fx.FxPure[V]
type AbortHn[V, E any] = func(fx.Fx[fx.And[AbortFn[V, E], fx.Nil], V]) fx.FxPure[Result[V, E]]
type AbortAb[V, E any] = fx.And[AbortFn[V, E], fx.Nil]
type AbortFx[V, E, U any] = fx.Fx[AbortAb[V, E], U]

func Abort[V, E any](e E) AbortFx[V, E, V] {
	return fx.Suspend[AbortFn[V, E]](e)
}

func Succeed[E, V any](v V) AbortFx[V, E, V] {
	return fx.Map(fx.Ctx[AbortAb[V, E]](), func(_ AbortAb[V, E]) V {
		return v
	})
}

func Handler[V, E any]() AbortHn[V, E] {
	return func(eff fx.Fx[fx.And[AbortFn[V, E], fx.Nil], V]) fx.FxPure[Result[V, E]] {
		var err Result[V, E]
		handler := fx.Handler(func(e E) fx.FxPure[V] {
			err = failure[V](e)
			return fx.Halt[fx.Nil, V]()
		})
		succeeded := fx.Map(handler(eff), success[V, E])
		failed := func() fx.FxPure[Result[V, E]] { return fx.Pure(&err) }
		return fx.Replace(succeeded, failed)
	}
}

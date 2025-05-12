package abort

import fx "github.com/vic/eff.go"

type Result[V, E any] func() (*V, *E)

func success[V, E any](v V) Result[V, E] { return func() (*V, *E) { return &v, nil } }
func failure[V, E any](e E) Result[V, E] { return func() (*V, *E) { return nil, &e } }

type AbortRq[V, E any] = func(E) fx.FxPure[V]
type AbortHn[V, E any] = func(AbortFx[V, E]) fx.Fx[fx.Nil, Result[V, E]]
type AbortAb[V, E any] = fx.And[AbortRq[V, E], fx.Nil]
type AbortFx[V, E any] = fx.Fx[AbortAb[V, E], V]

func Abort[V, E any](e E) AbortFx[V, E] {
	return fx.Suspend[AbortRq[V, E]](e)
}

func Succeed[E, V any](v V) AbortFx[V, E] {
	return fx.Map(fx.Ctx[AbortAb[V, E]](), func(_ AbortAb[V, E]) V {
		return v
	})
}

func Handler[V, E any]() AbortHn[V, E] {
	return func(eff AbortFx[V, E]) fx.FxPure[Result[V, E]] {
		var err Result[V, E]
		handler := fx.Handler(func(e E) fx.FxPure[V] {
			err = failure[V](e)
			return fx.Halt[fx.Nil, V]()
		})
		succeeded := fx.Map(handler(eff), success[V, E])
		failed := func() fx.FxPure[Result[V, E]] { return fx.Pure(&err) }
		return fx.Restart(succeeded, failed)
	}
}

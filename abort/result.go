package abort

import fx "github.com/vic/fx.go"

type ResultFn[V, E any] func(AbortFx[E, V]) fx.FxPure[Result[V, E]]
type ResultAb[V, E any] = fx.And[ResultFn[V, E], fx.Nil]
type ResultFx[V, E any] = fx.Fx[ResultAb[V, E], Result[V, E]]

type Result[V, E any] func() (*V, *E)

func success[V, E any](v V) Result[V, E] { return func() (*V, *E) { return &v, nil } }
func failure[V, E any](e E) Result[V, E] { return func() (*V, *E) { return nil, &e } }

func AbortResult[V, E any](eff AbortFx[E, V]) ResultFx[V, E] {
	return fx.Handle[ResultFn[V, E]](eff)
}

func AbortHandler[V, E any](e AbortFx[E, V]) fx.FxPure[Result[V, E]] {
	var err Result[V, E]
	var abortFn AbortFn[E] = func(e E) fx.FxNil {
		err = failure[V](e)
		return fx.Halt[fx.Nil, fx.Nil]()
	}
	succeeded := fx.Map(fx.ProvideLeft(e, abortFn), success[V, E])
	failed := func() fx.FxPure[Result[V, E]] { return fx.Pure(err) }
	return fx.Replace(failed)(succeeded)
}

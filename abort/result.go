package abort

import fx "github.com/vic/fx.go"

type Result[V, E any] func() (*V, *E)

func success[V, E any](v V) Result[V, E] { return func() (*V, *E) { return &v, nil } }
func failure[V, E any](e E) Result[V, E] { return func() (*V, *E) { return nil, &e } }

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

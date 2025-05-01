package abort

import "github.com/vic/eff.go"

type Result[E, V any] func() (*V, *E)

type AbortAb[E, V any] = eff.Ability[E, Result[E, V], eff.Nil]
type AbortEff[E, V any] = eff.Eff[AbortAb[E, V], Result[E, V]]

func Abort[E, V any](e E) AbortEff[E, V] {
	return eff.Request[AbortEff[E, V]](e)
}

func HandleAbort[E, V any](e eff.Eff[AbortAb[E, V], V]) eff.Eff[eff.Nil, Result[E, V]] {
	success := eff.Map(e, func(v V) Result[E, V] {
		return func() (*V, *E) { return &v, nil }
	})
	type H = eff.Handler[E, Result[E, V], eff.Nil]
	var handler H = func(err E, cont eff.Cont[eff.Nil, Result[E, V]]) eff.Eff[eff.Nil, Result[E, V]] {
		failure := Result[E, V](func() (*V, *E) { return nil, &err })
		return eff.Pure(&failure)
	}
	return eff.Provide(success, handler.Ability())
}

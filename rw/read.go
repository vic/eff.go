package rw

import "github.com/vic/eff.go"

type Reader[S any] func() *S

type ReadRq[S any] struct{}
type ReadRs[S any] = *S
type ReadAb[S any] = eff.Ability[ReadRq[S], ReadRs[S], eff.Nil]

func Read[S any]() eff.Eff[ReadAb[S], ReadRs[S]] {
	return eff.Request[eff.Eff[ReadAb[S], ReadRs[S]]](ReadRq[S]{})
}

func ReadHandler[S any](r Reader[S]) eff.Handler[ReadRq[S], ReadRs[S], eff.Nil] {
	return func(_ ReadRq[S], f eff.Cont[eff.Nil, ReadRs[S]]) eff.Eff[eff.Nil, ReadRs[S]] {
		return f(r())
	}
}

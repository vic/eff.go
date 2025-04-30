package rw

import "github.com/vic/eff.go"

type Writer[S any] func(*S)

type WriteRq[S any] = *S
type WriteRs[S any] struct{}
type WriteAb[S any] = eff.Ability[WriteRq[S], WriteRs[S], eff.Nil]

func Write[S any](v *S) eff.Eff[WriteAb[S], WriteRs[S]] {
	return eff.Request[eff.Eff[WriteAb[S], WriteRs[S]]](WriteRq[S](v))
}

func WriteHandler[S any](w Writer[S]) eff.Handler[WriteRq[S], WriteRs[S], eff.Nil] {
	return func(v WriteRq[S], f eff.Cont[eff.Nil, WriteRs[S]]) eff.Eff[eff.Nil, WriteRs[S]] {
		w(v)
		return f(WriteRs[S]{})
	}
}

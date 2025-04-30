package rw

import "github.com/vic/eff.go"

type Reader[S any] func() *S
type Writer[S any] func(*S)

type ReadRq[S any] struct{}
type ReadRs[S any] = *S
type ReadAb[S any] = eff.Ability[ReadRq[S], ReadRs[S], eff.Nil]
type ReadEff[S any] = eff.Eff[ReadAb[S], ReadRs[S]]

type WriteRq[S any] = *S
type WriteRs[S any] struct{}
type WriteAb[S any] = eff.Ability[WriteRq[S], WriteRs[S], eff.Nil]
type WriteEff[S any] = eff.Eff[WriteAb[S], WriteRs[S]]

func Read[S any]() ReadEff[S] {
	return eff.Request[ReadEff[S]](ReadRq[S]{})
}

func Write[S any](v *S) WriteEff[S] {
	return eff.Request[WriteEff[S]](WriteRq[S](v))
}

func ReadHandler[S any](r Reader[S]) eff.Handler[ReadRq[S], ReadRs[S], eff.Nil] {
	return func(_ ReadRq[S], f eff.Cont[eff.Nil, ReadRs[S]]) eff.Eff[eff.Nil, ReadRs[S]] {
		return f(r())
	}
}

func WriteHandler[S any](w Writer[S]) eff.Handler[WriteRq[S], WriteRs[S], eff.Nil] {
	return func(v WriteRq[S], f eff.Cont[eff.Nil, WriteRs[S]]) eff.Eff[eff.Nil, WriteRs[S]] {
		w(v)
		return f(WriteRs[S]{})
	}
}

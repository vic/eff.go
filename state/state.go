package state

import "github.com/vic/eff.go"

type Reader[S any] func() *S
type Writer[S any] func(*S)

type ReadRq[S any] struct{}
type ReadRs[S any] = *S
type ReadEff[S any] = eff.SusEff[ReadRq[S], ReadRs[S], eff.None]

type WriteRq[S any] = *S
type WriteRs[S any] struct{}
type WriteEff[S any] = eff.SusEff[WriteRq[S], WriteRs[S], eff.None]

func Read[S any]() ReadEff[S] {
	return eff.Suspend[ReadEff[S]](ReadRq[S]{})
}

func Write[S any](v *S) WriteEff[S] {
	return eff.Suspend[WriteEff[S]](WriteRq[S](v))
}

func ReadHandler[S any](r Reader[S]) eff.Handler[ReadRq[S], ReadRs[S], eff.None] {
	return func(_ ReadRq[S], f eff.Cont[eff.None, ReadRs[S]]) eff.Eff[eff.None, ReadRs[S]] {
		return f(r())
	}
}

func WriteHandler[S any](w Writer[S]) eff.Handler[WriteRq[S], WriteRs[S], eff.None] {
	return func(v WriteRq[S], f eff.Cont[eff.None, WriteRs[S]]) eff.Eff[eff.None, WriteRs[S]] {
		w(v)
		return f(WriteRs[S]{})
	}
}

package rw

import fx "github.com/vic/eff.go"

type Reader[S any] func() *S

type readRq[S any] = func(fx.Nil) fx.FxPure[*S]
type readHn[S any] = func(fx.Fx[fx.And[readRq[S], fx.Nil], *S]) fx.FxPure[*S]
type ReadAb[S any] = fx.And[readHn[S], fx.Nil]
type ReadFx[S any] = fx.Fx[ReadAb[S], *S]

func Read[S any]() ReadFx[S] {
	return fx.Request[readHn[S]](fx.PNil)
}

func ReadHandler[S any](r Reader[S]) readHn[S] {
	return fx.Handler(func(_ fx.Nil) fx.FxPure[*S] {
		v := r()
		return fx.Pure(&v)
	})
}

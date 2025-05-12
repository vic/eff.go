package rw

import (
	fx "github.com/vic/eff.go"
)

type Writer[S any] func(*S)

type writeRq[S any] = func(*S) fx.FxNil
type writeHn[S any] = func(fx.Fx[fx.And[writeRq[S], fx.Nil], fx.Nil]) fx.FxNil
type WriteAb[S any] = fx.And[writeHn[S], fx.Nil]
type WriteFx[S any] = fx.Fx[WriteAb[S], fx.Nil]

func Write[S any](v *S) WriteFx[S] {
	return fx.Request[writeHn[S]](v)
}

func WriteHandler[S any](w Writer[S]) writeHn[S] {
	return fx.Handler(func(s *S) fx.FxNil {
		w(s)
		return fx.PureNil
	})
}

package rw

import fx "github.com/vic/fx.go"

type Reader[T any] func() *T

type ReadFn[T any] = func(fx.Nil) fx.FxPure[*T]
type ReadHn[T any] = func(fx.Fx[fx.And[ReadFn[T], fx.Nil], *T]) fx.FxPure[*T]
type ReadAb[T any] = fx.And[ReadHn[T], fx.Nil]
type ReadFx[T, V any] = fx.Fx[ReadAb[T], V]

func Read[T any]() ReadFx[T, *T] {
	return fx.Request[ReadHn[T]](fx.PNil)
}

func ReadHandler[T any](r Reader[T]) ReadHn[T] {
	return fx.Handler(func(_ fx.Nil) fx.FxPure[*T] {
		v := r()
		return fx.Pure(v)
	})
}

package rw

import fx "github.com/vic/fx.go"

type Reader[T any] func() *T

type ReadFn[T any] func(fx.Nil) fx.FxPure[*T]
type ReadAb[T any] = fx.And[ReadFn[T], fx.Nil]
type ReadFx[T, V any] = fx.Fx[ReadAb[T], V]

func Read[T any]() ReadFx[T, *T] {
	return fx.Handle[ReadFn[T]](fx.PNil)
}

func ReadService[T any](r Reader[T]) ReadFn[T] {
	return func(_ fx.Nil) fx.FxPure[*T] { return fx.Pure(r()) }
}

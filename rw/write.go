package rw

import (
	fx "github.com/vic/fx.go"
)

type Writer[T any] func(*T)

type WriteFn[T any] = func(*T) fx.FxNil
type WriteHn[T any] = func(fx.Fx[fx.And[WriteFn[T], fx.Nil], fx.Nil]) fx.FxNil
type WriteAb[T any] = fx.And[WriteHn[T], fx.Nil]
type WriteFx[T, V any] = fx.Fx[WriteAb[T], V]

func Write[T any](v *T) WriteFx[T, fx.Nil] {
	return fx.Request[WriteHn[T]](v)
}

func WriteHandler[T any](w Writer[T]) WriteHn[T] {
	return fx.Handler(func(t *T) fx.FxNil {
		w(t)
		return fx.PureNil
	})
}

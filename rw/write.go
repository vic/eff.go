package rw

import (
	fx "github.com/vic/fx.go"
)

type Writer[T any] func(*T)

type WriteFn[T any] func(*T) fx.FxNil
type WriteAb[T any] = fx.And[WriteFn[T], fx.Nil]
type WriteFx[T, V any] = fx.Fx[WriteAb[T], V]

func Write[T any](v *T) WriteFx[T, fx.Nil] {
	return fx.Suspend[WriteFn[T]](v)
}

func WriteService[T any](w Writer[T]) WriteFn[T] {
	return func(t *T) fx.FxNil {
		w(t)
		return fx.PureNil
	}
}

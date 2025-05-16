package fx

type And[A, B any] func() (A, B)

func swap[A, B any](ab And[A, B]) And[B, A] {
	var ba And[B, A] = func() (B, A) {
		a, b := ab()
		return b, a
	}
	return ba
}

func left[A, B any](ab And[A, B]) A {
	a, _ := ab()
	return a
}

func right[A, B any](ab And[A, B]) B {
	_, b := ab()
	return b
}

func AndNil[S, V any](e Fx[S, V]) Fx[And[S, Nil], V] {
	return Then[And[S, Nil], V](left, Const)(e)
}

func AndContra[A, B, V any](e Fx[A, V], f func(B) A) Fx[B, V] {
	return ContraMap[V](f)(e)
}

func AndSwap[A, B, V any](e Fx[And[A, B], V]) Fx[And[B, A], V] {
	return Then[And[B, A], V](swap[B, A], Const)(e)
}

func AndFlat[A, B, V any](e Fx[A, Fx[B, V]]) Fx[And[A, B], V] {
	return FlatMap(e, identity)
}

func AndNest[A, B, V any](e Fx[And[A, B], V]) Fx[A, Fx[B, V]] {
	return Pending(func(a A) Fx[A, Fx[B, V]] { return Const[A](ProvideLeft(e, a)) })
}

func AndCollapse[A, V any](e Fx[And[A, A], V]) Fx[A, V] {
	return Pending(func(a A) Fx[A, V] { return ProvideLeft(e, a) })
}

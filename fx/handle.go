package fx

func Apply[F ~func(I) O, I, O any](i I) Fx[F, O] {
	return Map(Ctx[F](), func(f F) O { return f(i) })
}

func Suspend[A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[A, B], O] {
	return AndFlat(Apply[A](i))
}

func Handler[A ~func(I) Fx[B, O], B, I, O any](a A) func(Fx[And[A, B], O]) Fx[B, O] {
	return func(e Fx[And[A, B], O]) Fx[B, O] { return ProvideLeft(e, a) }
}

func Handle[F ~func(Fx[And[A, B], O]) Fx[B, O], A ~func(I) Fx[B, O], B, I, O any](i I) Fx[And[F, B], O] {
	return Suspend[F](Suspend[A](i))
}

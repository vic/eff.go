package fx

// Creates an stopped effect from a resume function
func Stop[S, V any](e func() Fx[S, V]) Fx[S, V] {
	return Fx[S, V]{res: e}
}

// Resume an effect if it was previously stopped.
func Resume[S, V any](e Fx[S, V]) Fx[S, V] {
	if e.res != nil {
		return e.res()
	}
	if e.imm != nil {
		return Const[S](e.imm())
	}
	return Pending(func(s S) Fx[S, V] { return Resume(e.sus(s)) })
}

// Replace with y if x is already Halted. Otherwise x continues.
func Replace[S, V any](y func() Fx[S, V]) func(Fx[S, V]) Fx[S, V] {
	return func(x Fx[S, V]) Fx[S, V] {
		if x.res != nil {
			return y()
		}
		if x.imm != nil {
			return Const[S](x.imm())
		}
		return Pending(func(s S) Fx[S, V] { return Replace(y)(x.sus(s)) })
	}
}

// An stopped effect that panics if resumed.
// Only useful with Replace.
//
// For example, an Abort effect halts since it has no possible value for V
// but then its Handler can Replace the halted effect with an Error value.
// See: abort/result.go
func Halt[S, V any]() Fx[S, V] {
	return Stop(func() Fx[S, V] {
		return Fx[S, V]{
			imm: func() V {
				panic("tried to Resume a halted effect. try using Replace instead")
			},
		}
	})
}

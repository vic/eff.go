# Effect Requirements

So far, we have seen that an effect `Fx[S, V]` can have `S` be `Nil` for effects that can be evaluated right away and non-`Nil` for those pending effects that still need to be provided some value.

In this chapter we will talk about Composite requirement types using the `And` type. And how different types of functions/effects can represent the very same computation. We also look at similarities with higher-order functions in functional-programming and how rotating or re-arranging effect requirements is indifferent in `Fx.go`. Finally, we show some the`And*` and `Provide*` combinators that can help you reshape your effect requirements.

## `And[A, B]` composed Requirement Types

Using the same "length of string" function from the previous chapters, we can describe it in different ways. 

```go
// This is an strange way of writing `func(string) int`.
// But this shape is used to understand the types bellow.
//
// Focus on what are the requirements needed to perform 
// a computation, more than the shape of the type.
//
// In particular, note that `Fx[Nil, V]` is like a `func() V`
func LenghtOfString(s string) func() int {
    return func() int { return len(s) }
}

type LenFn = func(string) Fx[Nil, int]
```

Note that all of the following types are equivalent, as they describe the very same requirements and result types:
- `func(string) int`
- `Fx[string, int]`
- `func(string) func() int`
- `func(string) Fx[Nil, int]`
- `Fx[string, Fx[Nil, int]]`
- `Fx[And[string, Nil], int]`
- `Fx[And[Nil, string], int]`

The last three examples represent nested effects and are equivalent to functions of arity > 1 or functions that return functions.

`And[A, B]` is the requirement for both `A` and `B` abilities. Notice on the last two examples, that they have their components swapped, however, its important to note that in `Fx.go`, _the *order* of the abilities on the requirement does not matter_ and can be freely swapped/joined/unjoined. More on this when we talk about `And*` combinators.

Also, note that `And[A, Nil]` is equivalent to just `A`. All of these types represent the same type of computation and an effect can be transformed to any of those types freely.



## `>1` arity functions as effects.

Suppose you have a function that multiplies an string length by n.

```go
func MulLen(s string, n int) int {
    return len(s) * n
}
```

`MulLen` can be described by the following types:

- `func(string, int) int`
- `func(string) func(int) int`
- `Fx[And[string, int], int]`
- `Fx[string, Fx[int, int]]`
- `Fx[int, Fx[string, int]]`
- `Fx[And[int, string], int]`

An important thing to note is that in `Fx`, the *requirements are identified by their type* and not by their name, so they can be freely re-arranged or provided in any order. Note that `And[X, X]` is equivalent to just a single `X` requirement, and that `And[And[X, Y], And[Y, X]]` is also equivalent to `And[X, Y]`.


## `And*` Combinators.

There are some functions (more will be added as they are found useful) that help you re-arrange `And`ed effect requirements:

```go
// Since `And[A, A]` is equivalent to just `A`.
// Used to collapse Nil requirements just before evaluation.
func AndCollapse(Fx[And[A, A], V]) Fx[A, V]

// Ands S with Nil in the effect requirements
func AndNil(Fx[S, V]) Fx[And[S, Nil], V]

// Swaps A and B. Note: this has no impact on how computation is
// performed, since requirements can be freely re-arranged.
func AndSwap(Fx[And[A, B], V]) Fx[And[B, A], V]


// FlatMaps the inner effect into the outer by 
// Anding their requirements.
func AndJoin(Fx[A, Fx[B, V]]) Fx[And[A, B], V]

// Reverse of Join
func AndDisjoin(Fx[And[A, B], V]) Fx[A, Fx[B, V]]

```

## `Provide*` Combinators.

These functions are used to provide requirements into effects. The result is another effect (no computation has been performed yet) but with less requirements.

```go
// Discharges the single S requirement.
func Provide(Fx[S, V], *S) Fx[Nil, V]


// Discharges the requirement of A by providing it.
func ProvideLeft(Fx[And[A, B], V], *A) Fx[B, V]

// Discharges the requirement of B by providing it.
func ProvideRight(Fx[And[A, B], V], *B) Fx[A, V]

// Discharges both A and B
func ProvideBoth(Fx[And[A, B], V], *A, *B) Fx[Nil, V]



// Provides A, the `CAAR` of the requirements list.
func ProvideA(Fx[And[And[A, C], And[B, D]], V], *A) Fx[And[C, And[B, D]], V]

// Provides B, the `CADR` of the requirements list.
func ProvideB(Fx[And[And[A, C], And[B, D]], V], *B) Fx[And[And[A, C], D], V]

// Provides A and B, the `CAAR` and `CADR` of the requirements list.
func ProvideAB(Fx[And[And[A, C], And[B, D]], V], *A, *B) Fx[And[C, D], V]
```


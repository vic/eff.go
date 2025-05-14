# Handlers

A _Handler_ is an effect transformer function of type `func(Fx[R, U]) Fx[S, V]`. Handlers are free to change the effect requirements, they tipically reduce the requirement set, but they could also introduce new requirements. Also the result value can be changed or be the same.

## Handling an effect

Lets re-write our previous "length of string" function as a Handler.

```go
type LenFn = func(string) fx.Fx[fx.Nil, int]

// Code is type annotated for clarity
var lenFx fx.Fx[fx.And[LenFn, fx.Nil], int] = fx.Suspend[LenFn]("hello")

// type is not needed but just added for clarity.
type LenHn = func(fx.Fx[fx.And[LenFn, fx.Nil], int]) fx.Fx[fx.Nil, int]

var handler LenHn = fx.Handler(func(s string) fx.Fx[fx.Nil, int] {
    truth := len(s)
    return fx.Pure(&truth)
})

// apply the handler to lenFx
var x *int = fx.Eval(handler(lenFx))
assert(*x == 5)
```

As you might guess, `fx.Handler` is just a wrapper for `ProvideLeft(Fx[And[Fn, S], O], *Fn) Fx[S, O]` where `Fn = func(I) Fx[S, O]`, an request-effect function.


## Requesting Handlers (effect-transformers) from the environment.

Of course, you can also request that a particular effect transformer (Handler) be present as a requirement of some computation. In this way the handler is provided only once but can be applied anywhere it is needed inside the program.

```go
// Same examples as above with some more types for clarity
type LenFn = func(string) fx.Fx[fx.Nil, int]
type LenFx = fx.Fx[fx.And[LenFn, fx.Nil], int]
type LenHn = func(LenFx) fx.Fx[fx.Nil, int]

var lenFx LenFx = fx.Suspend[LenFn]("hello")

// Request that an implementation of LenHn transformer 
// is available in the environment and it be applied to lenFx.
var lenAp Fx[And[LenHn, Nil], int] = fx.Apply[LenHn](lenFx)

var handler LenHn = fx.Handler(func(s string) fx.Fx[fx.Nil, int] {
    truth := len(s)
    return fx.Pure(&truth)
})

// Now instead of applying the handler directly to each effect it
// must handle, we provide it into the environment.
var provided Fx[Nil, int] = fx.ProvideLeft(lenAp, &handler)
val x *int = fx.Eval(provided)
assert(*x == 5)
```
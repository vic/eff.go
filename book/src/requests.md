# Effect Requests

Another way of creating effects in `Fx.go` is via an *effect-request* function.

A function of type `func(I) Fx[S, O]` is said to take an _effect-request_ `I` and produce an *suspended* effect `Fx[S, O]`.

For example, the function `func(HttpReq) Fx[HttpService, HttpRes]` states that given an `HttpReq` request you can obtain an `HttpRes` response, *provided* that an `HttpService` is available.

Using the "length of string" example of the previous chapters, we can use it to model an effect request:

```go
type LenFn = func(string) fx.Fx[fx.Nil, int]

// Code is type annotated for clarity
var lenFx fx.Fx[fx.And[LenFn, fx.Nil], int] = fx.Suspend[LenFn]("hello")
```

Note that `Suspend` takes the _type_ of a request-effect function and the request value for it. And yields a *suspended* effect of type `Fx[And[LenFn, Nil], int]`. The computation is said to be *suspended* because it knows not what particular implementation of `LenFn` should be used, and because of this, `LenFn` is part of the requirements, along with `Nil` the ability requirement on the result of `LenFn`.

Different implementations of `LenFn` can be provided to the `lenFx` effect.

```go
var bad LenFn = func(_ string) fx.Fx[fx.Nil, int] {
    lies := 42
    return fx.Pure(&lies)
}
var good LenFn = func(s string) fx.Fx[fx.Nil, int] {
    truth := len(s)
    return fx.Pure(&truth)
}

var x *int = fx.Eval(fx.ProvideLeft(lenFx, &bad))
assert(*x == 42)

var y *int = fx.Eval(fx.ProvideLeft(lenFx, &good))
assert(*y == 5)
```

Notice that by delaying which implementation of `LenFn` is used, the `lenFx` program description includes the effect-request `"hello"` and knows the general form of its response `Fx[Nil, int]`, but knows nothing about which particular interpretation of `LenFn` will be used.
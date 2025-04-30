# eff.go

An Algebraic Effects System for Golang.

> Experimental. API is still evolving and new effects will be added as they are discovered to be useful to the golang ecosystem.

### How are Algebraic Effects useful?

Algebraic Effects are useful because they allow programs to
be expressed not only in terms of what kind of value they can
compute but also on what possible side-effects or external resources will such a computation require.

By using Effect Handlers, the interpretation of how an effect is performed is independent of the program description. In this way, a single program description can be interpreted by a *test-handler* that could, for example, mock request to external services, and a *prod-handler* that could actually perform such requests.

If you want to read more about different language implementations and theory behind effects, read the [effects-bibliography](https://github.com/yallop/effects-bibliography).

`eff.go` is inspired by the following two implementations, and uses a similar notion of the _Handler_, _Ability_, and _Effect_ concepts:

- [Unison Abilities](https://www.unison-lang.org/docs/language-reference/abilities-and-ability-handlers/)
- [Kyo (Scala)](https://github.com/getkyo/kyo/)

# Tour

This section will try to introduce you to the concepts of
_Effects_, _Abilities_ and _Handlers_ in `eff.go` and how they can be used to describe your Golang programs.

No knowledge or previous experience with other effect sytems
is expected, and we will try to explain things inductively, by
working out from simple concepts to more interesting ones.


## Effects, Abilities and Handlers.

An _Effect_ (`Eff[S, V]` can be read _`V` provided `S`_) is the description of a computation of type `V` provided that the ability requirements `S` are present, so that the computation of `V` can be performed.

`S` is said to be the _Ability_ (or set of Abilities) that are needed for computing `V`. Abilities describe the external resources that would be needed as well as side-effects that are possible while computing `V`.

A `Handler` for the `S` Ability, provides a particular interpretation of what `S` means. It is the Handler that actually decides how to perform world-modifying side-effects.
It is possible and quite common to have different interpretations (or Handlers) of a single Ability, for example, for test and production runs.


> An `Effect` is just the recipe of a program (`V`).
It describes the `Abilities` (`S`) that are needed for producing `V`, but an Effect by itself does nothing. It is until a particular `Handler` of `S` is provided that the computation of `V` is actually executed.


### The `eff.Eff[S, V]` type

An effect `Eff[S, V]` can be one of two possible values:

- An *Immediate value: `V`*. That is, the value `V` has already been computed, and there's nothing to be done to determine it. No external resources nor side-effects are needed for it.

  Immediate values are created using the function `Pure[V](v *V)` that takes a pointer to an already existing value `V`.

  The pointer of an immediate value can be retrieved using the function `Eval[V](eff Eff[Nil, V]) *V` which takes an effect with the `Nil` ability requirement.

- A *Suspended value: `V` provided `S`*. That is, the computation of `V` is still pending, and `S` is needed for it to be completed.

  The most basic suspended computation is one you are already familiar with: *A Function*. For example:

  The function that computes an string length:

  ```go
  func StringLength(s string) int {
    return len(s)
  }
  ```

  Can be expressed as an effect of type `Eff[string, int]`. That is, in order to compute the `int` value you need first to be provided with an `string` value.

  ```go
  import ( . "github.com/vic/eff.go" )

  // Our first effect program from a traditional function.
  var eff1 Eff[string, int] = Func(StringLength)

  var requirement string := "hello"
  // Notice that the effect requirement is discharged
  var eff2 Eff[Nil, int] = Provide(eff1, &requirement)

  // Only effects depending on Nil can be evaled.
  result := Eval(eff2)

  // Dereference the immediate value.
  *result == len("hello")
  ```

  Using `eff.Func(func (S) V)` you could lift a function into their effect type.

  However, suspended values are actually of interest when we are using them with Abilities and Handlers.

### The `eff.Ability` and `eff.Handler` types.

An `Ability` is the description of external services or side-effects that are needed for a computation.

For example, lets write a program that needs to perform http requests in order to complete. Such a program would create an http-request and expect an http-response from the web service it accesses. On an effects system like `eff.go`, we do not directly contact external services, we just express our need to perform such requests, and expect a `Handler` to actually decide how and when such requests should be performed (if any).


```go
package example

// Notice our program does not depend on any http library,
// just effects.
import ( . "github.com/vic/eff.go" )

// This type represents an http-request.
// For simplicty we use a single string: the URL
// of a GET request.
type HttpRq string

// This type represents of an http-response.
// For simplicity we use a single string: the response body
type HttpRs string

// The Http Ability, specifies that we will:
// - make HttpRq requests
// - expect HttpRs responses
// - and that this ability requires `Nil` other abilities.
type HttpAb = Ability[HttpRq, HttpRs, Nil]

// Produce an effect of fetching the given URL
func Get(url string) Eff[HttpAb, HttpRs] {
    // eff.Request takes an HttpRq and produces
    // a suspended (delayed) computation of HttpRs
    return Request[Eff[HttpAb, HttpRs]](HttpRq(url))
}

func BodyLength(r HttpRs) int {
    return len(r)
}

// Computes the length of response from http://example.org
func Program() Eff[HttpAb, int] {
    e := Get("http://example.org")
    return Map(e, BodyLength)
}
```

When we invoke `Program()` only the `Eff[HttpAb, int]` recipe
is created, but no request is made to the external world, since we don't even include any http library.

In order to actually produce requests, we need a `Handler` for the `HttpAb` ability. The handler is the interpreter of `HttpRq` requests and produces actual `HttpRs` responses.

Lets assume we are creating tests and we wont actually interact with the network, since our program needs not to know where we get the responses from.

```go
package example

import (
    "testing"
    . "github.com/vic/eff.go"
)

// A test handler for the HttpAb ability. mocks responses.
// Notice this type has the same type parameters as HttpAb.
type HttpHandler = Handler[HttpRq, HttpRs, Nil]

func HttpTestHandler() HttpHandler {
    // A handler is nothing more than a function that takes:
    // - the ability request (HttpRq)
    // - a continuation that will create an Eff[Nil, HttpRs].
    //   In this example, the Nil requirement means that no
    //   other abilities are needed to compute the HttpRs response.
    return func(rq HttpRq, cont Cont[Nil, HttpRs]) Eff[Nil, HttpRs] {
        // ignore request and produce a fixed response.
        return cont(HttpRs("hello"))
    }
}

func TestProgram(t *testing.T) {
    // Explicit types are shown only for clarity
    var program Eff[HttpAb, int] = Program()
    var handler HttpHandler = HttpTestHandler()
    var ability *HttpAb = handler.Ability()
    var handled Eff[Nil, int] = Provide(program, ability)
    var result *int = Eval(handled)
    if *result != len("hello") {
        t.Error("unexpected result")
    }
}

```

### Combining Effects.

> TODO: ContraMap, Map, MapM, FlatMap

### Providing requirements.

> TODO: Provide, ProvideLeft, ProvideRight, ProvideBoth, Rotate

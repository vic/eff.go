# Introduction

An Algebraic Effects System for Golang.

<div class="warning">
Fx.go is currently experimental. 

API surface is *very much* in flux and evolving.

New effects will be added as they are discovered to be useful in the golang ecosystem.
</div>


### How are Algebraic Effects useful?

Algebraic Effects are useful because they allow programs to
be expressed not only in terms of what kind of value they can
compute but also on what possible side-effects or external resources will such a computation require.

By using Effect Handlers, the interpretation of how an effect is performed is independent of the program description. This means that a single program description can be interpreted in different ways. For example, using a *test-handler* that mocks request to external services, or using a *live-handler* that actually performs such requests.

If you want to read more about different language implementations and theory behind effects, read the [effects-bibliography](https://github.com/yallop/effects-bibliography).

`Fx.go` is inspired by the following two implementations, and uses a similar notion of the _Handler_, _Ability_, and _Effect_ concepts:

- [Unison Abilities](https://www.unison-lang.org/docs/language-reference/abilities-and-ability-handlers/)
- [Kyo (Scala3)](https://github.com/getkyo/kyo/) - special thanks to [@fbrasisil](https://x.com/fbrasisil), Kyo's author who kindly provided a minimal kyo-core that helped [me](https://x.com/oeiuwq) understand algebraic effect systems and inspired this library.


However, `Fx.go` has a different surface API since we are trying to provide the best dev-experience for Golang programmers.
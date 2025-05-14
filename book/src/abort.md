# Abort

`Abort[V, E](E)` aborts a computation of `V` with `E`.
The abort Handler transforms these effects into `Result[V, E]`.

- Implementation [abort.go](https://github.com/vic/fx.go/blob/main/abort/abort.go)
- Tests [abort_test.go](https://github.com/vic/fx.go/blob/main/abort/abort_test.go)

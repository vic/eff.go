package fx

type Nil pnil
type pnil struct{}

var PNil Nil = Nil(pnil{})
var PureNil FxNil = Pure(PNil)

type FxNil = FxPure[Nil]

type FxPure[V any] = Fx[Nil, V]

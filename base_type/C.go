package base_type

type Number interface {
	Increase(n *Int) *Int
	Reduce(n *Number) *Number
	Multiply(n *Number) *Number
	Divide(n *Number) *Number
}

type C[K any, V any] struct {
	v V
}

func (c *C[K, V]) Get() V {
	return c.v
}

func (c *C[K, V]) Set(k K, v V) {
	c.v = v
}

package base_type

import "strconv"

type Float struct {
	c *C[any, float64]
}

func NewFloat(i float64) *Float {
	return &Float{c: &C[any, float64]{v: i}}
}

func (i *Float) Increase(i2 *Float) *Float {
	i3 := i.c.Get() + i2.c.Get()
	return NewFloat(i3)
}

func (i *Float) Reduce(i2 *Float) *Float {
	i3 := i.c.Get() - i2.c.Get()
	return NewFloat(i3)
}

func (i *Float) Multiply(i2 *Float) *Float {
	i3 := i.c.Get() * i2.c.Get()
	return NewFloat(i3)
}

func (i *Float) Divide(i2 *Float) *Float {
	i3 := i.c.Get() / i2.c.Get()
	return NewFloat(i3)
}

func (i *Float) ToSting(prec int) string {
	return strconv.FormatFloat(i.c.Get(), 'f', prec, 64)
}

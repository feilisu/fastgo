package base_type

import "strconv"

type Int struct {
	c *C[any, int]
}

func NewInt(i int) *Int {
	return &Int{c: &C[any, int]{v: i}}
}

func (i *Int) Increase(i2 *Int) *Int {
	i3 := i.c.Get() + i2.c.Get()
	return NewInt(i3)
}

func (i *Int) Reduce(i2 *Int) *Int {
	i3 := i.c.Get() - i2.c.Get()
	return NewInt(i3)
}

func (i *Int) Multiply(i2 *Int) *Int {
	i3 := i.c.Get() * i2.c.Get()
	return NewInt(i3)
}

func (i *Int) Divide(i2 *Int) *Int {
	i3 := i.c.Get() / i2.c.Get()
	return NewInt(i3)
}

func (i *Int) Remainder(i2 *Int) *Int {
	i3 := i.c.Get() % i2.c.Get()
	return NewInt(i3)
}

func (i *Int) ToSting() string {
	return strconv.Itoa(i.c.Get())
}

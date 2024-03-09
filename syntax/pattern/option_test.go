package pattern

import "testing"

type ComplicateStruct struct {
	// 必传字段
	field1 string
	// 可选字段，可传可不传
	field2 string
	field3 string
}

type ComplicateStructOption func(c *ComplicateStruct)

func WithField2(field2 string) ComplicateStructOption {
	return func(c *ComplicateStruct) {
		c.field2 = field2
	}
}

func WithField3(field string) ComplicateStructOption {
	return func(c *ComplicateStruct) {
		c.field3 = field
	}
}

func NewComplicateStruct(field1 string,
	opts ...ComplicateStructOption) *ComplicateStruct {
	res := &ComplicateStruct{
		field1: field1,
		field2: "这是我的默认值",
		field3: "这还是我的默认值",
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func TestOption(t *testing.T) {
	c := NewComplicateStruct("这是必传",
		WithField2("Field2自定义的值"))
	t.Log(c.field2)
}

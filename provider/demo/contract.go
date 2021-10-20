package demo

const Key = "yogo:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}

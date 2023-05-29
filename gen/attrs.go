package gen

type Attr struct {
	name  string
	value any
}

func (a *Attr) Name() string {
	return a.name
}

func (a *Attr) Value() any {
	return a.value
}

func NewAttr(name string, value any) *Attr {
	return &Attr{name, value}
}

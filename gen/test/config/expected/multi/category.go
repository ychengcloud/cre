package multi

type Category struct {
	Id   int32
	Name string
}

func (c *Category) Ops(name string) []string {
	switch name {
	case "id":
		return []string{}
	case "name":
		return []string{}
	default:
		return []string{}
	}
}

func (c *Category) GetName() string {
	return "category"
}

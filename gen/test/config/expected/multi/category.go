package multi

type Category struct {
	Id int32 `json:"id"`

	Name string `json:"name,omitempty"`
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

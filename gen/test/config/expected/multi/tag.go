package multi

type Tag struct {
	Id int32 `json:"id"`

	Name string `json:"name,omitempty"`
}

func (t *Tag) Ops(name string) []string {
	switch name {
	case "id":
		return []string{}
	case "name":
		return []string{}
	default:
		return []string{}
	}
}

func (t *Tag) GetName() string {
	return "tag"
}

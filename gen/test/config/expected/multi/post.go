package multi

type Post struct {
	Id int32 `json:"id"`

	Name       string      `json:"name,omitempty"`
	Author     *User       `json:"author,omitempty"`
	Categories []*Category `json:"categories"`
	Tags       []*Tag      `json:"tags"`
}

func (p *Post) Ops(name string) []string {
	switch name {
	case "id":
		return []string{
			"Eq",
			"In",
		}
	case "name":
		return []string{}
	case "author":
		return []string{}
	case "categories":
		return []string{}
	case "tags":
		return []string{}
	default:
		return []string{}
	}
}

func (p *Post) GetName() string {
	return "post"
}

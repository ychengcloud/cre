package multi

type Post struct {
	Id         int32
	Name       string
	Author     *User
	Categories []*Category
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
	default:
		return []string{}
	}
}

func (p *Post) GetName() string {
	return "post"
}

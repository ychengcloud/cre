package multi

type User struct {
	Id    int32
	Name  string
	Posts []*Post
}

func (u *User) Ops(name string) []string {
	switch name {
	case "id":
		return []string{}
	case "name":
		return []string{
			"Eq",
			"In",
		}
	case "posts":
		return []string{}
	default:
		return []string{}
	}
}

func (u *User) GetName() string {
	return "user"
}

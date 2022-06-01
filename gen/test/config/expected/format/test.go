package single

type Category struct {
}

func (c *Category) GetName() string {
	return "category"
}

type Post struct {
}

func (p *Post) GetName() string {
	return "post"
}

type User struct {
}

func (u *User) GetName() string {
	return "user"
}

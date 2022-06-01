package single

import (
	"github.com/ychengcloud/cre/gen/test/config/expected/path"
)

type Category struct {
	Id   int32
	Name string

	p path.Test
}

func (c *Category) GetName() string {
	return "category"
}

type Post struct {
	Id         int32
	Name       string
	Author     *User
	Categories []*Category

	p path.Test
}

func (p *Post) GetName() string {
	return "post"
}

type User struct {
	Id    int32
	Name  string
	Posts []*Post

	p path.Test
}

func (u *User) GetName() string {
	return "user"
}

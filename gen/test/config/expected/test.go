package single

import (
	"github.com/ychengcloud/cre/gen/test/config/expected/path"
)

type Category struct {
	Id int32 `json:"id"`

	Name string `json:"name,omitempty"`

	p path.Test
}

func (c *Category) GetName() string {
	return "category"
}

type Tag struct {
	Id int32 `json:"id"`

	Name string `json:"name,omitempty"`

	p path.Test
}

func (t *Tag) GetName() string {
	return "tag"
}

type Post struct {
	Id int32 `json:"id"`

	Name       string      `json:"name,omitempty"`
	Author     *User       `json:"author,omitempty"`
	Categories []*Category `json:"categories"`
	Tags       []*Tag      `json:"tags"`

	p path.Test
}

func (p *Post) GetName() string {
	return "post"
}

type User struct {
	Id int32 `json:"id"`

	Name  string  `json:"name,omitempty"`
	Posts []*Post `json:"posts"`

	p path.Test
}

func (u *User) GetName() string {
	return "user"
}

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	files := []string{
		"test.go",
		"multi/category.go",
		"multi/post.go",
		"multi/user.go",
		"path/test.go",
		"format/test.go",
		"format/test.proto",
	}
	r := require.New(t)

	for _, file := range files {
		r.FileExists(filepath.Join("actual", file))
		actual, err := os.ReadFile(filepath.Join("actual", file))
		r.NoError(err)
		r.NotEmpty(actual)

		expected, err := os.ReadFile(filepath.Join("expected", file))
		r.NoError(err)
		r.NotEmpty(expected)
		r.Equal(string(expected), string(actual))

	}

	// 解析代码源文件，获取常量和注释之间的关系
	// fset := token.NewFileSet()
	// f, err := parser.ParseFile(fset, "gen/multi/user.go", nil, 0)
	// r.NoError(err)
	// r.NotNil(f)

	// r.Equal(true, test.Exist(f, "User", "int32", "ID"))
	// r.Equal(true, test.Exist(f, "User", "Post", "Posts"))
	// r.Equal(false, test.Exist(f, "User", "", "noexist"))

	// fpostset := token.NewFileSet()
	// fpost, err := parser.ParseFile(fpostset, "gen/multi/post.go", nil, 0)
	// r.NoError(err)
	// r.NotNil(f)

	// r.Equal(true, test.Exist(fpost, "Post", "int32", "ID"))
	// r.Equal(true, test.Exist(fpost, "Post", "User", "Author"))
	// r.Equal(false, test.Exist(fpost, "Post", "", "noexist"))
}

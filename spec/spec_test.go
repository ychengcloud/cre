package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOps(t *testing.T) {
	id := &Field{
		Name:       "id",
		Type:       &IntegerType{Name: "int", Size: 32},
		PrimaryKey: true,
	}

	name := &Field{
		Name: "name",
		Type: &StringType{Name: "char", Size: 64},
	}

	ops := defaultOps(id.Type, id.Optional)
	r := require.New(t)
	r.Equal(len(NumericOps), len(ops))
	r.Equal(Eq, ops[Eq-1])
	r.Equal(Neq, ops[Neq-1])
	r.Equal(In, ops[In-1])
	r.Equal(NotIn, ops[NotIn-1])
	r.Equal(Gt, ops[Gt-1])
	r.Equal(Gte, ops[Gte-1])
	r.Equal(Lt, ops[Lt-1])
	r.Equal(Lte, ops[Lte-1])

	ops = defaultOps(name.Type, name.Optional)
	r.Equal(len(StringOps), len(ops))
}

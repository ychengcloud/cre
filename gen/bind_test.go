package gen

import (
	"testing"

	"github.com/ychengcloud/cre"
	"gotest.tools/assert"
)

func TestRebind(t *testing.T) {
	b := &Binder{Dialect: cre.Postgres}
	assert.Equal(t, b.Rebind("SELECT * FROM user WHERE id = ?"), "SELECT * FROM user WHERE id = $1")
	assert.Equal(t, b.Rebind("SELECT * FROM user WHERE id = ? AND name = ?"), "SELECT * FROM user WHERE id = $1 AND name = $2")
	assert.Equal(t, b.Rebind("SELECT * FROM user WHERE id = ? AND name = ? AND age = ?"), "SELECT * FROM user WHERE id = $1 AND name = $2 AND age = $3")
	assert.Equal(t, b.Rebind("Insert into user (id, name, age) values (?, ?, ?)"), "Insert into user (id, name, age) values ($1, $2, $3)")
}

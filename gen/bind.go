package gen

import (
	"bytes"
	"strconv"

	"github.com/ychengcloud/cre"
)

type Binder struct {
	Dialect string
}

func (b *Binder) Rebind(query string) string {
	if b.Dialect == cre.MySQL {
		return query
	}

	qb := make([]byte, 0, len(query))
	rb := bytes.NewBuffer(qb)

	i := 1
	for _, q := range query {
		if q == '?' {
			rb.WriteRune('$')
			rb.WriteString(strconv.Itoa(i))
			i++
		} else {
			rb.WriteRune(q)
		}
	}
	return rb.String()
}

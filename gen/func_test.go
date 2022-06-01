package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncs(t *testing.T) {
	require.Equal(t, "Id", pascal("id"))
	require.Equal(t, "id", camel("id"))
	require.Equal(t, "ids", plural("id"))
	require.Equal(t, "id", singular("ids"))

	require.Equal(t, "f", receiver([]string{"fmt"}, "fmt"))
	require.Equal(t, "fm", receiver([]string{"fmt", "f"}, "fmt"))
	require.Equal(t, "_fmt", receiver([]string{"fmt", "f", "fm"}, "fmt"))
}

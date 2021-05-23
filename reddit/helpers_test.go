package reddit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOrderedMaxSet(t *testing.T) {
	set := NewOrderedMaxSet(1)
	set.Add("foo")
	set.Add("bar")
	println(len(set.keys))
	require.Equal(t, set.Len(), 1)
}

func TestOrderedMaxSetCollision(t *testing.T) {
	set := NewOrderedMaxSet(2)
	set.Add("foo")
	set.Add("foo")

	require.Equal(t, set.Len(), 1)
}

func TestOrderedMaxSet_Delete(t *testing.T) {
	set := NewOrderedMaxSet(1)
	set.Add("foo")

	require.Equal(t, set.Len(), 1)

	set.Delete("foo")
	require.Equal(t, set.Len(), 0)
	require.False(t, set.Exists("foo"))
}

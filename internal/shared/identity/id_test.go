package identity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfrr/randomtalk/internal/shared/identity"
)

func TestID_String(t *testing.T) {
	id := identity.ID("test-id")
	expected := "test-id"
	assert.Equal(t, expected, id.String())
}

func TestID_IsEmpty(t *testing.T) {
	id1 := identity.ID("")
	id2 := identity.ID("test-id")

	assert.True(t, id1.IsEmpty())
	assert.False(t, id2.IsEmpty())
}

func TestID_Equals(t *testing.T) {
	id1 := identity.ID("test-id")
	id2 := identity.ID("test-id")
	id3 := identity.ID("another-id")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}

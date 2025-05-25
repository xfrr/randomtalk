package identity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfrr/randomtalk/internal/shared/identity"
)

func TestNewUUID(t *testing.T) {
	id := identity.NewUUID()
	assert.NotEmpty(t, id)
	assert.NotEqual(t, id, identity.ID(""))
	assert.NotEqual(t, id, identity.NewUUID())
}

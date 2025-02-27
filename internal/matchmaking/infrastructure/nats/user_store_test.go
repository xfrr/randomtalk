package matchnats_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	matchnats "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/nats"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

func TestUserStore_AddUser(t *testing.T) {
	ctx := context.Background()
	js := setupJetStream(t)
	store, err := matchnats.NewUserStore(ctx, js)
	require.NoError(t, err)

	user := matchdomain.NewUser("user-id-1", 25, gender.GenderUnspecified, matchmaking.DefaultPreferences())
	err = store.AddUser(ctx, user)
	require.NoError(t, err)

	// Verify user was added
	users, err := store.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, user.ID(), users[0].ID())
}

func TestUserStore_GetAll(t *testing.T) {
	ctx := context.Background()
	js := setupJetStream(t)
	store, err := matchnats.NewUserStore(ctx, js)
	require.NoError(t, err)

	user1 := matchdomain.NewUser("user-id-1", 25, gender.GenderUnspecified, matchmaking.DefaultPreferences())
	user2 := matchdomain.NewUser("user-id-2", 30, gender.GenderUnspecified, matchmaking.DefaultPreferences())

	err = store.AddUser(ctx, user1)
	require.NoError(t, err)
	err = store.AddUser(ctx, user2)
	require.NoError(t, err)

	users, err := store.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserStore_RemoveUsers(t *testing.T) {
	ctx := context.Background()
	js := setupJetStream(t)
	store, err := matchnats.NewUserStore(ctx, js)
	require.NoError(t, err)

	user := matchdomain.NewUser("user-id-1", 25, gender.GenderUnspecified, matchmaking.DefaultPreferences())
	err = store.AddUser(ctx, user)
	require.NoError(t, err)

	err = store.RemoveUsers(ctx, user.ID())
	require.NoError(t, err)

	users, err := store.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, users)
}

func setupJetStream(t *testing.T) jetstream.JetStream {
	t.Helper()

	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)

	js, err := jetstream.New(nc)
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = js.DeleteKeyValue(ctx, "randomtalk_matchamking_user_store")
		require.NoError(t, err)

		nc.Close()
	})

	return js
}

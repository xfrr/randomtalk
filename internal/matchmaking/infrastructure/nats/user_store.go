package matchnats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

var _ matchdomain.UserStore = &UserStore{}

// UserStore is the nats implementation of the UserStore interface
type UserStore struct {
	js jetstream.JetStream
	kv jetstream.KeyValue
}

// AddUser implements matchdomain.UserStore.
func (u *UserStore) AddUser(ctx context.Context, user matchdomain.User) error {
	// 1. serialize the user
	body, err := matchdomain.MarshalUser(&user)
	if err != nil {
		return err
	}

	// 2. upsert the user
	_, err = u.kv.Put(ctx, user.ID(), body)
	if err != nil {
		return err
	}

	return err
}

// GetAll implements matchdomain.UserStore.
func (u *UserStore) GetAll(ctx context.Context) ([]*matchdomain.User, error) {
	keys, err := u.kv.Keys(ctx, jetstream.IgnoreDeletes())
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrNoKeysFound):
			return []*matchdomain.User{}, nil
		default:
			return nil, fmt.Errorf("failed to get keys: %w", err)
		}
	}

	users := make([]*matchdomain.User, 0, len(keys))
	for _, key := range keys {
		userFound, err := u.kv.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		if userFound.Operation() == jetstream.KeyValueDelete {
			// ignore deleted keys
			continue
		}

		user, err := matchdomain.UnmarshalUser(userFound.Value())
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// RemoveUsers implements matchdomain.UserStore.
func (u *UserStore) RemoveUsers(ctx context.Context, userID ...string) error {
	revisions := make([]uint64, 0, len(userID))
	for _, id := range userID {
		user, err := u.kv.Get(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, jetstream.ErrKeyNotFound):
				return matchdomain.ErrUserNotFound
			default:
				return err
			}
		}

		revisions = append(revisions, user.Revision())
	}

	for i, id := range userID {
		rev := revisions[i]
		err := u.kv.Delete(ctx, id, jetstream.LastRevision(rev))
		if err != nil {
			return err
		}
	}

	return nil
}

// NewUserStore creates a new UserStore
func NewUserStore(ctx context.Context, js jetstream.JetStream) (*UserStore, error) {
	// create a new kv store
	kvstore, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "randomtalk_matchamking_user_store",
		History: 1,
		TTL:     1 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	return &UserStore{
		js: js,
		kv: kvstore,
	}, nil
}

package matchdomain

import "github.com/google/uuid"

type MatchID string

func (id MatchID) String() string {
	return string(id)
}

func GenerateID() MatchID {
	return MatchID(uuid.New().String())
}

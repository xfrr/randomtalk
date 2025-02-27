package identity

// ID represents an unique string identifier for an entity.
type ID string

// String returns the ID as a string.
func (id ID) String() string {
	return string(id)
}

// Equals checks if two IDs are equal.
func (id ID) Equals(other ID) bool {
	return id == other
}

// IsEmpty checks if the ID is empty.
func (id ID) IsEmpty() bool {
	return id == ""
}

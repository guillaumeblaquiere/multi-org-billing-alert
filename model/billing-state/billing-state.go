package billing_state

type State int

const (
	Active State = iota
	Deleted
)

func (s State) String() string {
	return [...]string{"active", "deleted"}[s]
}

func (s *State) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

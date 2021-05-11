package actions

// Action type
type Action string

func (e Action) String() string {
	return string(e)
}

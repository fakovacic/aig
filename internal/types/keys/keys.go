package keys

type key string

const (
	// UUID const
	UUID key = "UUID"
)

func (c key) String() string {
	return string(c)
}

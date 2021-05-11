package methods

type Method string

func (e Method) String() string {
	return string(e)
}

const (
	// Get const
	Get Method = "GET"
	// Post const
	Post Method = "POST"
	// Patch const
	Patch Method = "PATCH"
	// Delete const
	Delete Method = "DELETE"
)

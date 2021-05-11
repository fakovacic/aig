package endpoints

import (
	"fmt"
	"strings"
)

// Endpoint type
type Endpoint string

func (e Endpoint) String() string {
	return string(e)
}

// Format endpoint to url
func (e Endpoint) Format(a ...interface{}) string {
	return fmt.Sprintf(strings.ReplaceAll(e.String(), ":id", "%s"), a...)
}

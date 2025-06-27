package state

import (
	"fmt"
	"strings"
)

type State string

type sessionData map[string]any

func (d sessionData) String() string {
	var b strings.Builder
	for k, v := range d {
		b.WriteString(fmt.Sprintf("%s %v\n", k, v))
	}
	return b.String()
}

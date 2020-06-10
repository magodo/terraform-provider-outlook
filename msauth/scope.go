package msauth

import (
	"sort"
	"strings"
)

type scope struct {
	scopes []string
}

func NewScope(scopes ...string) scope {
	sort.Strings(scopes)
	return scope{scopes}
}

func (s scope) String() string {
	return strings.Join(s.scopes, " ")
}

func (s scope) IsEmpty() bool {
	return len(s.scopes) == 0
}

package containerutils

import (
	"path"
	"strings"
)

// Attribute is a string that has specific methods - syntactic sugar
type Attribute string

// Attributes is a slice of Attribute ([]string)
type Attributes []Attribute

const (
	mountDir   = "tmp"
	routerPort = "80:80"
)

func joinMountDir(filepath Attribute) Attribute {
	f := func(r rune) bool {
		return (r == '/' || r == ' ' || r == '\\')
	}
	dir := strings.TrimRightFunc(mountDir, f)
	p := path.Join(dir, string(filepath))
	return Attribute(p)
}

// And creates a new string slice with the given Attribute arguments
// appended
func (m Attributes) And(m1 ...Attribute) Attributes {
	return append(m, m1...)
}

// AndString creates a new string slice with the given string arguments
// appended (note that Attribute is a typedef for Go's native string type)
func (m *Attributes) AndString(m1 ...string) Attributes {
	m2 := make([]Attribute, len(m1)+len(*m))
	l1 := copy(m2, (*m))
	for i, _m := range m1 {
		m2[l1+i] = Attribute(_m)
	}
	*m = m2
	return m2
}

// Equals returns the given slice m as a slice containing only m
// - syntactic sugar
func (m Attribute) Equals(m2 string) bool {
	return m.String() == m2
}

// Contains checks if attribute array contains string m2
func (m Attributes) Contains(m2 string) bool {
	for _, _m := range m {
		if _m.Equals(m2) {
			return true
		}
	}
	return false
}

// ToSlice returns the given slice m as a slice containing only m
// - syntactic sugar
func (m Attribute) ToSlice() Attributes {
	return []Attribute{m}
}

// String recovers the underlying string of the Attribute
func (m Attribute) String() string {
	return string(m)
}

// Remove an element from the attributes
func (m *Attributes) Remove(value string) {
	if m == nil || len(*m) == 0 {
		return
	}
	outArray := make([]Attribute, len(*m)-1)
	offset := 0
	for i := 0; i < len(*m); i++ {
		if (*m)[i].String() == value {
			offset++
			if len(*m)+offset >= len(*m) {
				break
			}
		}
		outArray[i] = (*m)[i+offset]
	}
	*m = Attributes(outArray)
}

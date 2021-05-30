package inplaceenvsubst

// Simple implementation of a string set
type StringSet map[string]bool

func NewStringSet(values ...string) StringSet {

	set := map[string]bool{}
	for _, v := range values {
		set[v] = true
	}
	return set
}

// returns true if value exists in the set
// false otherwise
func (s StringSet) Contains(value string) bool {
	_, exists := s[value]
	return exists
}

func (s StringSet) IsEmpty() bool{
	return s == nil || len(s) == 0
}

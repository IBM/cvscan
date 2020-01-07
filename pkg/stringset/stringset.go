package stringset

import (
	"encoding/json"
	"fmt"
	"sort"
)

type set = map[string]struct{}

// StringSet is a hash set of strings.
type StringSet struct {
	set set
}

// New instantiates a new empty StringSet and returns a reference to it.
func New(keys ...string) *StringSet {
	ss := &StringSet{
		set: set{},
	}
	ss.Add(keys...)
	return ss
}

// Add inserts a key into the set.
func (ss *StringSet) Add(keys ...string) {
	if ss == nil {
		return
	}
	for _, k := range keys {
		ss.set[k] = struct{}{}
	}
}

// Delete removes a key from the set.
func (ss *StringSet) Delete(keys ...string) {
	if ss == nil {
		return
	}
	for _, k := range keys {
		delete(ss.set, k)
	}
}

// Has returns whether or not the set contains the given string.
func (ss *StringSet) Has(key string) bool {
	if ss == nil {
		return false
	}
	_, exists := ss.set[key]
	return exists
}

// Len returns the number of elements in the set.
func (ss *StringSet) Len() int {
	if ss == nil {
		return 0
	}
	return len(ss.set)
}

// AsSlice returns an unordered slice of the set.
func (ss *StringSet) AsSlice() []string {
	slc := []string{}
	if ss == nil {
		return slc
	}
	for s := range ss.set {
		slc = append(slc, s)
	}
	return slc
}

// Set returns the raw map data.
func (ss *StringSet) Set() set { return ss.set }

func (ss *StringSet) String() string {
	if ss == nil {
		return ""
	}
	keys := ss.AsSlice()
	sort.Strings(keys)
	return fmt.Sprintf("%q", keys)
}

func (ss *StringSet) MarshalJSON() ([]byte, error) {
	slice := ss.AsSlice()
	sort.Strings(slice)
	return json.Marshal(slice)
}

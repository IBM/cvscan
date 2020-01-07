package stringset

import (
	"fmt"
	"testing"
)

func TestHas(t *testing.T) {
	empty := New()
	nonempty := New(
		"key1",
		"key2",
		"my other key",
	)

	cases := []struct {
		name     string
		ss       *StringSet
		k        string
		expected bool
	}{
		{
			name:     "nil set empty key",
			ss:       nil,
			k:        "",
			expected: false,
		},
		{
			name:     "empty set empty key",
			ss:       empty,
			k:        "",
			expected: false,
		},
		{
			name:     "empty set nonempty key",
			ss:       empty,
			k:        "key1",
			expected: false,
		},
		{
			name:     "nonempty set empty key",
			ss:       nonempty,
			k:        "",
			expected: false,
		},
		{
			name:     "nonempty set key does not exist",
			ss:       nonempty,
			k:        "does not exist",
			expected: false,
		},
		{
			name:     "nonempty set key exists",
			ss:       nonempty,
			k:        "key1",
			expected: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.ss.Has(test.k)
			if actual != test.expected {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	(*StringSet)(nil).Delete("a")

	ss := New("s")
	ss.Delete("s")
	if ss.Has("s") {
		t.Error("set still Has Deleted key")
	}

	ss = New("a", "b", "c")
	ss.Delete("b")
	if ss.Has("b") {
		t.Error("Deleted b but it still exists")
	}

	if !ss.Has("a") {
		t.Error("deleted b but a does not exist")
	}

	if !ss.Has("c") {
		t.Error("deleted b but c does not exist")
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		ss       *StringSet
		expected int
	}{
		{
			ss:       nil,
			expected: 0,
		},
		{
			ss:       New(),
			expected: 0,
		},
		{
			ss:       New("a"),
			expected: 1,
		},
		{
			ss:       New("a", "b"),
			expected: 2,
		},
		{
			ss:       New("a", "b", "c"),
			expected: 3,
		},
		{
			ss:       New("a", "b", "c", "d", "e", "f", "g", "h", "i", "j"),
			expected: 10,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.expected), func(t *testing.T) {
			if test.ss.Len() != test.expected {
				t.Errorf("expected length %d, got %d", test.expected, test.ss.Len())
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	var ss *StringSet
	for i := 0; i < b.N; i++ {
		ss = New()
	}
	_ = ss
}

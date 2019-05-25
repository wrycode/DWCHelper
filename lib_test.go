package main

import (

	"testing"
	"encoding/json"
)

func TestRemove (t *testing.T) {

	var removeTests = []struct {
		terms []string // original slice
		term string // term to remove
		out []string // returned slice
	}{
		{[]string{"a","b","c"}, "a", []string{"b","c"}},
		{[]string{"a","b","c"}, "b", []string{"a","c"}},
		{[]string{"a","b","c"}, "c", []string{"a","b"}},
		{[]string{"a","b","b","c"}, "b", []string{"a","c"}},
		{[]string{"a"}, "a", []string{}},
		{[]string{}, "a", []string{}},
		{[]string{"a","b","c"}, "x", []string{"a","b","c"}},
	}
	
	for _, tt := range removeTests {
		result, _ := json.Marshal(Remove(tt.terms, tt.term))
		expected, _ := json.Marshal(tt.out)
		if string(result) != string(expected) {
			t.Errorf("Remove(%v, %v): expected %v, got %v", tt.terms, tt.term, string(expected), string(result))
		}
	}
}

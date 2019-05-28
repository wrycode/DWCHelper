package main

import (

	"testing"
	"encoding/json"
	"strings"
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

func TestRename (t *testing.T) {
	var renameTests = []struct {
		terms []string // original slice
		oldTerm string // term to rename
		newTerm string // new name
		out []string // returned slice
	}{
		{[]string{"a","b","c"}, "a","b", []string{"b","b","c"}},
		{[]string{"a","b","c"}, "x","b", []string{"a","b","c"}},
		{[]string{"a","b","b","c"}, "b","x", []string{"a","x","x","c"}},
		{[]string{"a"}, "a", "", []string{""}},
		{[]string{}, "a","x", []string{}},
	}

	for _, tt := range renameTests {
		result, _ := json.Marshal(Rename(tt.terms,tt.oldTerm,tt.newTerm))
		expected, _ := json.Marshal(tt.out)
		if string(result) != string(expected) {
			t.Errorf("Rename(%v, %v, %v): expected %v, got %v", tt.terms, tt.oldTerm,tt.newTerm, string(expected), string(result))
		}
	}
	
}

func TestInputNumber(t *testing.T) {
	var inputNumberTests = []struct {
		first int // lower limit
		second int //upper limit
		choice string // user input
		result int // number returned
	}{
		{0, 5, "5", 5},
		{0, 5, "0", 0},
		{0, 5, "a\n4", 4},
		{0, 5, "asaf\nasf\n4", 4},
		{0, 5, "6\n0", 0},
		{0, 5, "-400\n0", 0},
	}

	for _, test := range inputNumberTests {
		reader := strings.NewReader(test.choice)
		result := inputNumber(test.first,test.second,reader)
		if result != test.result {
			t.Errorf("inputNumber(%v,%v) with user input: %v \n got %v, expected %v",test.first,test.second,test.choice,result,test.result)
		}
	}
}

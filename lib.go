package main

// Index returns the first index of the target string t, or -1 if no
// match is found
func Index(terms []string, s string) int {
	for i, v := range terms {
		if v == s {
			return i
		}
	}
	return -1
}

// Include returns true if the target string t is in the slice
func Include(terms []string, term string) bool {
	return Index(terms, term) >= 0
}

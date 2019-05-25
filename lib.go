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

// Replace returns a []string with oldTerm replaced by newTerm
// func Replace( terms []string, oldTerm, newTerm) []string {
// }

// Remove returns a []string with all instances of term removed, or
// unchanged if the term isn't in the slice
func Remove(terms []string, term string) []string {
	var result []string
	for _, t := range terms {
		if t != term {
			result = append(result, t)
		}
	}
	if result != nil {
		return result
	}
	return []string{}
}

package main

import (
	"io"
	"os"
	"strconv"
	"bufio"
	"fmt"
	"time"
	"strings"
	"github.com/fatih/camelcase"
)

const termURL string = "https://raw.githubusercontent.com/tdwg/dwc/master/dist/simple_dwc_horizontal.csv"

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

// Rename renames all occurrences of oldTerm to newTerm in terms
func Rename(terms []string, oldTerm, newTerm string) []string {
	for i, t := range terms {
		if t == oldTerm {
			terms[i] = newTerm
		}
	}
	return terms
}

// inputNumber returns an int between the given upper and lower
// limits, taken from the io.Reader argument (such as os.Stdin). If
// the input is invalid, it returns 0
func inputNumber (first int, second int, r io.Reader) int {
	fmt.Printf("Your choice? (%v to %v): ",first,second)
	b := bufio.NewScanner(r)
	for b.Scan() {
		n, err := strconv.Atoi(b.Text())
		if err == nil {
			if first <= n && n <= second {
				fmt.Println()
				return n
			}
		}
		fmt.Printf("Please enter a valid number between %v and %v and hit Enter:\n", first, second)
	}
	if err := b.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println()
	return 0
}

// printStringSlice prints a []string, surrounding each string with
// quotes
func printStringSlice(terms []string) {
	for i, v := range terms {
		fmt.Printf("\"%v\" ",v)
		if i % 3 == 0 {
			fmt.Println()
		}
	}
}

// Prompt prints out the given string, and asks for user confirmation
// to continue is ask is set to true
func Prompt(ask bool, s string) {
	b := bufio.NewScanner(strings.NewReader(s))
		time.Sleep(25 * time.Millisecond)


	for b.Scan() {
		time.Sleep(25 * time.Millisecond)
		fmt.Println(b.Text())
	}
	if ask {
		r := bufio.NewReader(os.Stdin)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Press Enter to continue...")
		for {
			_, err := r.Peek(1) 
			if err == nil {
				return
			}
		}
	}
}

// PrintHLine draws one or more horizontal lines
func PrintHLine(i int) {
	for n := 0; n < i; n = n +1 {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}
}


// notAllSame returns true if not every element of a string slice is the same
func notAllSame(s []string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != s[0] {
			return true
		}
	}
	return false
}

// stringIsVariation is a loose matching function to determine if a
// string could be a variation of a term in camelCase
func stringIsVariation(s string, t string) bool {

	words := camelcase.Split(t)
	
	switch {
	case s == strings.ToLower(strings.Join(words, " ")):
		return true
	case s == strings.Title(strings.Join(words, " ")):
		return true
	}

	stringWords := strings.Split(s, " ")

	for _, x := range stringWords {
		for _, y := range words {
			if strings.ToLower(x) == strings.ToLower(y) {
				return true
			}
		}
	}

	return false
}

// inputTerm gets a new term from the user
func inputTerm(message string, r io.Reader) string {
	fmt.Print(message)
	b := bufio.NewScanner(r)
	b.Scan()
	fmt.Println()
	return b.Text()
}

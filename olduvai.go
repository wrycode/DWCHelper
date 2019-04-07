// Olduvai interactively imports data from a CSV file (given as a
//command-line argument) into our standardized database format.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Check for filename argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: olduvai <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Get some data to work with
	data, err := ioutil.ReadFile(filename)
	if err != nil { // Exit if file doesn't exist
		fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	variables := strings.Split(lines[0], ",") // List of variables
	lines = lines[1:]                         // all other lines of the data
	lines = lines[:len(lines) - 1] // removes the final line, which is an empty string

	var dbase1 database

	dbase1.data = make(map[string]column)

	// Build database from CSV (TODO: make more efficient with pointers?)
	for i, variable := range variables {
		var temp column
		var values []string
		for _, line := range lines {
			values = append(values, strings.Split(line, ",")[i])
		}
		temp.values = values
		dbase1.data[variable] = temp
	}

	// TODO: Various checks on the dataset 

	// Import DWC aliases and definitions

	// Run through each variable and build new database with user input

}

// holds all of the variables and their data 
type database struct {
	data map[string]column // maps variable names to their full column
	// aliases 
}
// holds one "column" of the database; aka one variable and all of its
// values(in order)
type column struct {
	varType, alias, definition string // other metadata
	values []string // values are stored in strings
}

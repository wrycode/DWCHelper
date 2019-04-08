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
	dbase1 := importData(filename)
	//	fmt.Println(checkDatabase(dbase1))
	for i, value := range dbase1.variables {
		fmt.Print(value, " ")
		if i % 5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
//	fmt.Print(dbase1.variables)
	
	// TODO: Various checks on the dataset 

	// Import DWC aliases and definitions

	// Run through each variable and build new database with user input

}
// importData imports a CSV file. It takes a filename as an argument and returns a database
func importData(filename string) database {
	// Get some data to work with
	data, err := ioutil.ReadFile(filename)
	if err != nil { // Exit if file doesn't exist
		fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	vars := strings.Split(lines[0], ",") // List of variables
	lines = lines[1:]                         // all other lines of the data
	lines = lines[:len(lines) - 1] // removes the final line, which is an empty string

	//initialize database and data field
	var dbase database 
	dbase.data = make(map[string]column)

	// build database
	numOfVars := len(vars)
	for i, variable := range vars {
		var temp column
		var values []string
		for _, line := range lines {
			x := strings.Split(line, ",")
       			if len(x) != numOfVars { // checks that each specimen has one value for each variable
				fmt.Println("Error while importing data.")
				fmt.Printf("Double check the following line, it has %v values and it should have %v\n",
					len(x), numOfVars)
				fmt.Println(line)
				os.Exit(1)
			}
			values = append(values, x[i])
		}
		temp.values = values // set the column field for the current variable
		dbase.data[variable] = temp // write the column to the database
	}

	dbase.variables = vars // add the ordered list of variables (dbase.data is a map, and therefore unordered)
	return dbase
}

// holds all of the variables and their data 
type database struct {
	data map[string]column // maps variable names to their full column
	variables []string
	// aliases 
}

// holds one "column" of the database; aka all of the values for one
// variable (in order)
type column struct {
	varType, alias, definition string // other metadata
	values []string // values are stored in strings
}

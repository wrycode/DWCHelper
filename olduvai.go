// Olduvai interactively imports data from a CSV file (given as a
//command-line argument) into our standardized database format.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

func main() {
	// Check for filename argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: olduvai <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	originalDB := importData(filename) // original dataset, which we won't change

	// copy originalDB over to a fresh database for modification
	var alteredDB database
	alteredDB.data = make(map[string]column)
	for k, v := range originalDB.data {
		alteredDB.data[k] = v
	}
	alteredDB.variables = originalDB.variables

	// set "hasDifferentValues" variable for each column
	for s, v := range alteredDB.data {
		v.hasDifferentValues = notAllSame(v.values)
		alteredDB.data[s] = v
	}
	
	// populate the "type" for each column
	for s, v := range alteredDB.data {
			alteredDB.data[s] = inferType(v)
	}

	// print info about "useful" variables
	for s, v := range alteredDB.data {
		if v.hasDifferentValues {
			fmt.Println("Variable name: ", s)
			fmt.Println("=============================================================")
			fmt.Println("Values: ", v.values[50:55])
			fmt.Println("Inferred type: ", v.varType)
			fmt.Println("v.hasDifferentValues: ", v.hasDifferentValues)
			fmt.Println()
		}
	}
	
	
	// Import DWC aliases and definitions

	// Run through each variable and build new database with user input
	fmt.Println()
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

// notAllSame returns true if not every element of a string slice is the same
func notAllSame(s []string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != s[0] {
			return true
		}
	}
	return false
}

// inferTypes fills up the varType variable for a column based on
// built in parsing, also remove surrounding quotes from values
func inferType(c column) column {
	c.varType = "string"
//	p := c.hasDifferentValues // Only print output about the inferences if it's a useful variable
	for _, value := range c.values {
	//	if p { p := (i == 0 | 50 | 100) } // Print info for only some of the values

		// Strip quotes
		if strings.HasSuffix(value, "\"") && strings.HasPrefix(value, "\"") {
			// if p {
			// fmt.Printf("Stripping quotes from value %v , ",value)
			// }
			value = strings.Trim(value, "\"")
		}

		// Attempt number conversions
		_, err1 := strconv.ParseFloat(value, 64)
		_, err2 := strconv.Atoi(value)

		if err1 == nil {
			// if p {
			// 	fmt.Printf(" setting numeric/floating point type, ")
			// }
			c.varType = "numeric/floating point"
		}
		if err2 == nil {
			// if p {
			// 	fmt.Printf(" setting integer type ")
			// }
			c.varType = "integer"
		}
	}
//	fmt.Printf("varType: %v")
	return c
}

// holds all of the variables and their data 
type database struct {
	data map[string]column // maps variable names to their full column
	variables []string // original ordered list of variables
}

// holds one "column" of the database; aka all of the values for one
// variable (in order)
type column struct {
	values []string // values are stored in strings
	varType, alias, definition string // other metadata
	hasDifferentValues bool // whether the values change for each specimen (good indicator that it's a useful variable)
}

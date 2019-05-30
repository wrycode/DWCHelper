/* DWCHelper interactively formats CSV data for Simple Darwin Core
(https://dwc.tdwg.org/simple/) compatibility.

Run DWCHelper with two command-line arguments, the first being the
input file and the second being the output file.  */
package main

import (        
	"encoding/csv"
	"fmt"
	"runtime"
	"os"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	//	"bufio"

)


const referenceURL string = "http://rs.tdwg.org/dwc/terms/"
const aliasURL = "https://git.sr.ht/~wrycode/DWCHelper/blob/master/aliases.csv"

func main() {
	// Check for filename argument
	if len(os.Args) != 3 {
		fmt.Println("Usage: DWCHelper <input-filename.csv> <output-filename.csv>")
		os.Exit(1)
	}

	// Import database from file given as first command-line argument
	db := importDB(os.Args[1])

	// check for .settings file, if it exists, apply the saved
	// settings.  Otherwise, run the helper functions

	f, err := os.Open(os.Args[1] + ".settings")
	defer f.Close()
	if err == nil {
		Prompt(false,`Using settings from previous run. To run with
clean options and redo the import process, please` +
	"delete " + os.Args[1] + ".settings "+ " and re-run DWCHelper...")

		r := csv.NewReader(f)
		r.LazyQuotes = true
		
		// .settings file is not "square"
		r.FieldsPerRecord = -1

		// remove terms
		termsToRemove, err := r.Read()
		if err != nil {
			fmt.Println("Cannot read CSV data for terms to remove in the settings file:", err.Error())
			os.Exit(1)
		}
		
		for _, val := range termsToRemove {
			db = removeTerm(val, db)
		}

		// rename the rest
		rows, err := r.ReadAll()
		if err != nil {
			fmt.Println("Cannot read CSV data for aliases in the settings file:", err.Error())
			os.Exit(1)
		}

		for _, row := range rows {
			alias := row[0]
			newTerm := row[1]
			db = renameTerm(alias, newTerm, db)
		}
	} else {
		// CSV writer to save the settings in the file
		settingsFile, err := os.Create(os.Args[1] + ".settings")
		var w *csv.Writer
		if err != nil {
			fmt.Printf("Cannot save settings to '%s': %s\n", os.Args[1] + ".settings", err.Error())
			fmt.Println("Proceeding without saving your conversion settings...")
			w = csv.NewWriter(ioutil.Discard)
		} else {
			w = csv.NewWriter(settingsFile)
		}

		// fix windows line endings
		if runtime.GOOS == "windows" {
			w.UseCRLF = true
		}

		// remove terms
		termsToRemove := removeHelper(db)
		for _, val := range termsToRemove {
			db = removeTerm(val, db)
		}

		// save removed terms
		w.Write(termsToRemove)

		// rename terms
		rows := renameHelper(db)
		for _, row := range rows {
			alias := row[0]
			DWCTerm := row[1]
			db = renameTerm(alias, DWCTerm, db)
		}

		// save renamed terms
		w.WriteAll(rows)
		w.Flush()
	}

	// Export database to file given as second command-line argument
	exportDB(os.Args[2], db)

}

// importDB imports a CSV file. It takes a filename as an argument and returns a database
func importDB(filename string) database {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// TODO do I need to close the reader?
	r := csv.NewReader(f)
	r.LazyQuotes = true
	rows, err := r.ReadAll()
	if err != nil {
		fmt.Println("Cannot read CSV data:", err.Error())
		os.Exit(1)
	}

	// Initialize database
	var db database
	db.data = make(map[string]column)
	// Ordered list of terms
	db.terms = rows[0]
	// Fill in columns
	for i, term := range db.terms {
		var temp column
		var values []string
		for _, row := range rows[1:] {
			values = append(values, row[i])
		}
		temp.values = values
		db.data[term] = temp
	}
	return db
}

// removeTerm removes a given term from the database's list of terms
func removeTerm(term string, db database) database {
	if Include(db.terms, term) {
		fmt.Println("Removing",term)
		db.terms = Remove(db.terms, term)
	}
	return db
}

// removeHelper is the interactive helper function that returns a list
// of terms to be removed
func removeHelper(db database) []string {
	var termsToRemove []string

	for _, term := range db.terms {
		if !notAllSame(db.data[term].values) {
			termsToRemove = append(termsToRemove, term)
		}
	}
	PrintHLine(1)
	
	Prompt(false,`First we will clean up your list of terms. 
For the following terms , the values are either empty (no data), or the value is the same 
for every specimen:`)

	PrintHLine(1)

	//asff
	fmt.Println()
	printStringSlice(termsToRemove)
	fmt.Println()
	
	PrintHLine(1)
	Prompt(false, `Would you like to delete them?
0: no, don't delete any terms
1: yes, delete all of the above terms
2: delete some terms (let me choose)`)
	PrintHLine(1)

	//	switch n := inputNumber(0,2, os.Stdin); n {
	switch 1 {
	case 0:
		return []string{}
	case 2:
		choices := make([]bool, len(termsToRemove))
		PrintHLine(1)

		Prompt(false,`Which terms would you like to remove?
1 through ` + strconv.Itoa(len(termsToRemove)) + ": select a term" + `
-1: done entering terms
0: show which terms are currently selected for removal`)


	done := false
		PrintHLine(1)

		for done == false {
			switch n := inputNumber(-1, len(termsToRemove), os.Stdin); n {
			case 0:
				for i, v := range termsToRemove {
					fmt.Printf("%v: \"%v\" ",i+1,v)
					if choices[i] == true {
						fmt.Printf(" <===REMOVE ")
					}
					if i % 3 == 0 {
						fmt.Println()
					}
				}
				fmt.Println()
			case -1:
				done = true
			default: if choices[n- 1] == false {
				choices[n - 1] = true
			} else { choices[n - 1] = false}
				
			}
		}

		var chosenTerms []string
		for i, b := range choices {
			if b {
				chosenTerms = append(chosenTerms, termsToRemove[i])
			}
		}
		termsToRemove = chosenTerms
	}
	fmt.Println()
	return termsToRemove
}

// renameTerm renames a term in a given database (including the new
// mapping in the "data" field)
func renameTerm(oldName, newName string, db database) database {
	if Include(db.terms, oldName) {
		fmt.Printf("Renaming \"%v\" to \"%v\"\n",oldName,newName)
		db.data[newName] = db.data[oldName]
		db.terms = Rename(db.terms, oldName,newName)
	}
	return db
}

// renameHelper is the interactive helper function that returns a 2D
// array that maps terms to their new names
func renameHelper(db database) [][]string {
	var termsAndNewTerms [][]string
	var suggestions [][]string
	DWCTerms := pullDWCTerms()
	PrintHLine(1)
	Prompt(false,`These are the remaining terms. You can select a term by its 
number and rename it. Some terms have suggestions for names that 
have been  used by others. It may be helpful to refer to  the list
of terms at https://dwc.tdwg.org/terms/ while you do this.
(Pulling suggestions may take a few moments)`)
	PrintHLine(1)

	// generate suggestions for each term

	for i, term := range db.terms {
		// blank suggestions entry 
		suggestions = append(suggestions, []string{term})
		termsAndNewTerms = append(termsAndNewTerms, []string{term})
		// add DWCTerm if term may be a variation of it
		for _, DWCTerm := range DWCTerms {
			if stringIsVariation(term, DWCTerm) {
				suggestions[i] = append(suggestions[i], DWCTerm)
			}
		}

		// add alias from pre-generated aliases.csv file if
		// term may be a variation of it
		for _, row := range pullAliases() {
			for _, aliasTerm := range row[1:] {
				if stringIsVariation(term, aliasTerm) {
					suggestions[i] = append(suggestions[i], row[0])
				}
			}
		}
	}


	for i, row := range termsAndNewTerms {
		fmt.Printf(" %v: \"%v\" ",i+1, row[0])
		if len(termsAndNewTerms[i]) > 1 {
			fmt.Printf(" ======> %v",termsAndNewTerms[i][1])
		}
		if len(suggestions[i]) > 1 {
			fmt.Printf("(Suggestions: ")
			for _, suggestion := range suggestions[i][1:] {
				fmt.Printf("\"%v\" ",suggestion)
			}
			fmt.Printf(")")
		}
		fmt.Println()
	}
	done := false
	for done == false {
		fmt.Printf("-1: Done renaming | 0: list terms again | %v - %v: select term \n",1,len(termsAndNewTerms))
		switch n := inputNumber(-1, len(termsAndNewTerms), os.Stdin); n {
		case -1 : done = true
		case 0 :
			for i, row := range termsAndNewTerms {
				fmt.Printf(" %v: \"%v\" ",i+1, row[0])
				if len(termsAndNewTerms[i]) > 1 {
					fmt.Printf(" ======> %v",termsAndNewTerms[i][1])
				}
				if len(suggestions[i]) > 1 {
					fmt.Printf("(Suggestions: ")
					for _, suggestion := range suggestions[i][1:] {
						fmt.Printf("\"%v\" ",suggestion)
					}
					fmt.Printf(")")
				}
				fmt.Println()
			}
		default:
			if len(termsAndNewTerms[n-1]) > 1 {
				termsAndNewTerms[n - 1][1] = inputTerm("Please enter the new name for " + termsAndNewTerms[n - 1][0] + ": ", os.Stdin)
			} else {
				termsAndNewTerms[n - 1] = append(termsAndNewTerms[n - 1], 
					inputTerm("Please enter the new name for \"" + termsAndNewTerms[n - 1][0] + "\": ", os.Stdin))
			}
		}
	}
	
	var test [][]string
	for _, row := range termsAndNewTerms {
		if len(row) >= 2 {
			test = append(test, row)
		}
	}
	return test
}

// pullAliases pulls the alias database from the repository
func pullAliases() [][]string {
	resp, err := http.Get(aliasURL)
	if err != nil {
	 	fmt.Printf("Cannot pull aliases from upstream: %v\n", err.Error())
	 	fmt.Println("Skipping the pre-generated alias suggestions...")
	} else {
		contents, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("Error: cannot read contents of the alias file", err.Error())
			fmt.Println("Skipping the pre-generated alias suggestions...")
		} else {
			r := csv.NewReader(strings.NewReader(string(contents)))
			r.FieldsPerRecord = -1 // uneven fields allowed
			rows, err := r.ReadAll()
			if err != nil {
				fmt.Println("Cannot parse CSV in alias file:", err.Error())
				fmt.Println("Skipping the pre-generated alias suggestions...")
			} else {
				return rows
			}
		}
	}
	return [][]string{}
}

// pullDWCTerms grabs the current list of DWC Simple terms from their
// Github repository into a []string
func pullDWCTerms() []string {
	// default list of terms:
	termList := "type,modified,language,license,rightsHolder,accessRights,bibliographicCitation,references,institutionID,collectionID,datasetID,institutionCode,collectionCode,datasetName,ownerInstitutionCode,basisOfRecord,informationWithheld,dataGeneralizations,dynamicProperties,occurrenceID,catalogNumber,recordNumber,recordedBy,individualCount,organismQuantity,organismQuantityType,sex,lifeStage,reproductiveCondition,behavior,establishmentMeans,occurrenceStatus,preparations,disposition,associatedMedia,associatedReferences,associatedSequences,associatedTaxa,otherCatalogNumbers,occurrenceRemarks,organismID,organismName,organismScope,associatedOccurrences,associatedOrganisms,previousIdentifications,organismRemarks,materialSampleID,eventID,parentEventID,fieldNumber,eventDate,eventTime,startDayOfYear,endDayOfYear,year,month,day,verbatimEventDate,habitat,samplingProtocol,sampleSizeValue,sampleSizeUnit,samplingEffort,fieldNotes,eventRemarks,locationID,higherGeographyID,higherGeography,continent,waterBody,islandGroup,island,country,countryCode,stateProvince,county,municipality,locality,verbatimLocality,minimumElevationInMeters,maximumElevationInMeters,verbatimElevation,minimumDepthInMeters,maximumDepthInMeters,verbatimDepth,minimumDistanceAboveSurfaceInMeters,maximumDistanceAboveSurfaceInMeters,locationAccordingTo,locationRemarks,decimalLatitude,decimalLongitude,geodeticDatum,coordinateUncertaintyInMeters,coordinatePrecision,pointRadiusSpatialFit,verbatimCoordinates,verbatimLatitude,verbatimLongitude,verbatimCoordinateSystem,verbatimSRS,footprintWKT,footprintSRS,footprintSpatialFit,georeferencedBy,georeferencedDate,georeferenceProtocol,georeferenceSources,georeferenceVerificationStatus,georeferenceRemarks,geologicalContextID,earliestEonOrLowestEonothem,latestEonOrHighestEonothem,earliestEraOrLowestErathem,latestEraOrHighestErathem,earliestPeriodOrLowestSystem,latestPeriodOrHighestSystem,earliestEpochOrLowestSeries,latestEpochOrHighestSeries,earliestAgeOrLowestStage,latestAgeOrHighestStage,lowestBiostratigraphicZone,highestBiostratigraphicZone,lithostratigraphicTerms,group,formation,member,bed,identificationID,identificationQualifier,typeStatus,identifiedBy,dateIdentified,identificationReferences,identificationVerificationStatus,identificationRemarks,taxonID,scientificNameID,acceptedNameUsageID,parentNameUsageID,originalNameUsageID,nameAccordingToID,namePublishedInID,taxonConceptID,scientificName,acceptedNameUsage,parentNameUsage,originalNameUsage,nameAccordingTo,namePublishedIn,namePublishedInYear,higherClassification,kingdom,phylum,class,order,family,genus,subgenus,specificEpithet,infraspecificEpithet,taxonRank,verbatimTaxonRank,scientificNameAuthorship,vernacularName,nomenclaturalCode,taxonomicStatus,nomenclaturalStatus,taxonomy's"

	// Try to pull the csv termlist from online
	resp, err := http.Get(termURL)
	if err != nil {
		fmt.Printf("Cannot pull terms from Darwin Core repository: %s\n", err.Error())
		fmt.Println("Using the built-in terms instead...")
	} else {
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error: cannot read contents of termURL", err.Error())
		}
		termList = string(contents)
		defer resp.Body.Close()
	}

	// Use the CSV package to read the terms into a []string
	r := csv.NewReader(strings.NewReader(termList))
	terms, err := r.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return terms
}

// exportDB exports its database argument to the file at the filename argument
func exportDB(filename string, db database) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	// fix windows line endings
	if runtime.GOOS == "windows" {
		w.UseCRLF = true
	}
	w.Write(db.terms)                            // first line contains the terms in order
	for i := range db.data[db.terms[0]].values { // use the length of the first column as the number of rows
		var row []string
		for _, value := range db.terms { // for each term
			row = append(row, db.data[value].values[i]) // add the value of the term for the current row
		}
		if err := w.Write(row); err != nil { // write the row
			fmt.Println("error writing record to csv:", err)
		}
	}
	w.Flush()
}

// database holds all of the variables and their data
type database struct {
	data  map[string]column // maps terms to data
	terms []string          // ordered list of terms
}

// column holds one column of the database (not including the name of the term)
type column struct {
	values                     []string // values are stored in strings
	varType, alias, definition string   // other metadata
	hasDifferentValues         bool     // whether the values change for each specimen (good indicator that it's a useful variable)
}

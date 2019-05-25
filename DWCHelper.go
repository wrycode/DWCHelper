/* DWCHelper interactively formats CSV data for Simple Darwin Core
(https://dwc.tdwg.org/simple/) compatibility.

Run DWCHelper with two command-line arguments, the first being the
input file and the second being the output file.  */
package main

import (        
	"encoding/csv"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/fatih/camelcase"
)

const termURL string = "https://raw.githubusercontent.com/tdwg/dwc/master/dist/simple_dwc_horizontal.csv"
const referenceURL string = "http://rs.tdwg.org/dwc/terms/"
const aliasURL = "https://git.sr.ht/~wrycode/olduvai/blob/master/aliases.csv"

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
		fmt.Println("Using settings from previous run. To run with")
		fmt.Println("clean options and redo the import process, please")
		fmt.Println("delete", os.Args[1] + ".settings", " and re-run DWCHelper...")
		fmt.Println()

		r := csv.NewReader(f)
		r.LazyQuotes = true
		
		// .settings file is not "square"
		r.FieldsPerRecord = -1

		// remove terms
		termsToRemove, err := r.Read()
		if err != nil {
			fmt.Println("Cannot read CSV data in the settings file:", err.Error())
			os.Exit(1)
		}

		for _, val := range termsToRemove {
			db = removeTerm(val, db)
		}

		// TODO remember to flush various writers

		// rename the rest
		rows, err := r.ReadAll()
		if err != nil {
			fmt.Println("Cannot read CSV data in the settings file:", err.Error())
			os.Exit(1)
		}

		for _, row := range rows {
			alias := row[0]
			DWCTerm := row[1]
			db = renameTerm(alias, DWCTerm, db)
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

	fmt.Println(`First we will clean up your list of terms. 
The following terms are either empty (no data), or the value is the same for every 
specimen:`)

fmt.Println("Would you like to delete them?")
	fmt.Println("0: no, don't delete any terms")
	fmt.Println("1: yes, delete all of the above terms")
	fmt.Println("2: delete some terms (let me choose)")


	
	
	return []string{"Specimen number", "Comments"}
}

// renameTerm renames a term in a given database (including the new
// mapping in the "data" field)
func renameTerm(oldName, newName string, db database) database {
	fmt.Println("Renaming",oldName,"to",newName)
	db.data[newName] = db.data[oldName]
	
	return db
}

// renameHelper is the interactive helper function that returns a map
// of terms to their new names
func renameHelper(db database) [][]string {
	test := [][]string{ 
		{"Catalogue number", "catalogNumber"},
		{"Gnawing damage", "gnawingDamage"}}
	return test
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

// showReference returns the URL for a term on the Darwin Core website,
// which includes a definition, comments and examples
func showReference(term string) string {
	return referenceURL + "#" + term
}

// getAliases pulls the alias database from the repository or sets up
// a default one. It returns a map of alias strings to their possible term
func getAliases(terms []string) map[string]string { 
	//  map for storing possible aliases
	var aliases = make(map[string]string)
	
	// Try to pull the alias file from online
	resp, err := http.Get(aliasURL)
	if err != nil {
	 	fmt.Printf("Cannot pull aliases from upstream: %v\n", err.Error())
	 	fmt.Println("Defaulting to an automatically generated alias list...")
	} else { // Try to read the file
		contents, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
	 		fmt.Println("Error: cannot read contents of aliasURL", err.Error())
	 		fmt.Println("Defaulting to an automatically generated alias list...")
	
	 	} else { // Try to parse using csv
			r := csv.NewReader(strings.NewReader(string(contents)))
			r.FieldsPerRecord = -1 // uneven fields numbers allowed
			rows, err := r.ReadAll()
			if err != nil {
				fmt.Println("Cannot read CSV data from alias file:", err.Error())
				fmt.Println("Defaulting to an automatically generated alias list...")
			} else {

			// If there were no errors pulling the aliases
			// from the online CSV file, add obvious
			// variations of the alias to the aliases map
			for _, row := range rows {
				term := row[0]
				for _, entry := range row[1:] {
//					fmt.Println("running addAliases(",entry,"aliases",term)
					addAliases(entry, aliases, term)
					
				}
			}
			}
		}
			
	}

	// Generate default aliases from the DWC term names themselves
	for i := 1; i < 10; i++ {
		fmt.Println()
	}
	for _, term := range terms {
//		fmt.Println("running addAliases(",term,"aliases",term)

	addAliases(term, aliases, term)
}
	return aliases
}

// addAliases is a helper function for getAliases. It takes a word in
// camelCase and maps all obvious variations to the given aliases map
func addAliases(word string, aliases map[string]string, term string) {
words := camelcase.Split(word)
	aliases[strings.Join(words, "")] = term
	aliases[strings.Title(strings.Join(words, " "))] = term
	aliases[strings.ToLower(strings.Join(words, " "))] = term
//	fmt.Println(strings.Join(words, " "))
}

// showAliases is a temporary function (for debugging) that shows all
// aliases mapped to a term
func showAliases(term string, aliases map[string]string) {
	for key, value := range aliases {
		if value == term {
			fmt.Print(key,", ")
		}
	}
	for _, word := range camelcase.Split(term) {
		fmt.Print(word, ", ")
	}
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

// notAllSame returns true if not every element of a string slice is the same
func notAllSame(s []string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != s[0] {
			return true
		}
	}
	return false
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

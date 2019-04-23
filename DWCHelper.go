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
)

const termURL string = "https://raw.githubusercontent.com/tdwg/dwc/master/dist/simple_dwc_horizontal.csv"
const referenceURL string = "http://rs.tdwg.org/dwc/terms/"

func main() {
	// Check for filename argument
	if len(os.Args) != 3 {
		fmt.Println("Usage: DWCHelper <input-filename.csv> <output-filename.csv>")
		os.Exit(1)
	}

	// Import database from file given as first command-line argument
	db := importDB(os.Args[1])

	// TODO (number is how many hours I'm expecting for each task)
	// Type inference 4
	// Import DWC terms 2 DONE
	// Import DWC aliases WAITING on building alias file
	// DWC term inferences WAITING on importing terms and aliases
	// Detect unused terms, optionally remove them 2
	// Interactively check inferences, provide corrections 3? 4?

	// Export database to file given as second command-line argument
	exportDB(os.Args[2], db)
	fmt.Println(pullDWCTerms())
}

// importDB imports a CSV file. It takes a filename as an argument and returns a database
func importDB(filename string) database {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer f.Close()

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

// pullDWCTerms grabs the current list of DWC Simple terms from their
// Github repository into a []string
func pullDWCTerms() []string {
	// default list of terms:
	termList := "type,modified,language,license,rightsHolder,accessRights,bibliographicCitation,references,institutionID,collectionID,datasetID,institutionCode,collectionCode,datasetName,ownerInstitutionCode,basisOfRecord,informationWithheld,dataGeneralizations,dynamicProperties,occurrenceID,catalogNumber,recordNumber,recordedBy,individualCount,organismQuantity,organismQuantityType,sex,lifeStage,reproductiveCondition,behavior,establishmentMeans,occurrenceStatus,preparations,disposition,associatedMedia,associatedReferences,associatedSequences,associatedTaxa,otherCatalogNumbers,occurrenceRemarks,organismID,organismName,organismScope,associatedOccurrences,associatedOrganisms,previousIdentifications,organismRemarks,materialSampleID,eventID,parentEventID,fieldNumber,eventDate,eventTime,startDayOfYear,endDayOfYear,year,month,day,verbatimEventDate,habitat,samplingProtocol,sampleSizeValue,sampleSizeUnit,samplingEffort,fieldNotes,eventRemarks,locationID,higherGeographyID,higherGeography,continent,waterBody,islandGroup,island,country,countryCode,stateProvince,county,municipality,locality,verbatimLocality,minimumElevationInMeters,maximumElevationInMeters,verbatimElevation,minimumDepthInMeters,maximumDepthInMeters,verbatimDepth,minimumDistanceAboveSurfaceInMeters,maximumDistanceAboveSurfaceInMeters,locationAccordingTo,locationRemarks,decimalLatitude,decimalLongitude,geodeticDatum,coordinateUncertaintyInMeters,coordinatePrecision,pointRadiusSpatialFit,verbatimCoordinates,verbatimLatitude,verbatimLongitude,verbatimCoordinateSystem,verbatimSRS,footprintWKT,footprintSRS,footprintSpatialFit,georeferencedBy,georeferencedDate,georeferenceProtocol,georeferenceSources,georeferenceVerificationStatus,georeferenceRemarks,geologicalContextID,earliestEonOrLowestEonothem,latestEonOrHighestEonothem,earliestEraOrLowestErathem,latestEraOrHighestErathem,earliestPeriodOrLowestSystem,latestPeriodOrHighestSystem,earliestEpochOrLowestSeries,latestEpochOrHighestSeries,earliestAgeOrLowestStage,latestAgeOrHighestStage,lowestBiostratigraphicZone,highestBiostratigraphicZone,lithostratigraphicTerms,group,formation,member,bed,identificationID,identificationQualifier,typeStatus,identifiedBy,dateIdentified,identificationReferences,identificationVerificationStatus,identificationRemarks,taxonID,scientificNameID,acceptedNameUsageID,parentNameUsageID,originalNameUsageID,nameAccordingToID,namePublishedInID,taxonConceptID,scientificName,acceptedNameUsage,parentNameUsage,originalNameUsage,nameAccordingTo,namePublishedIn,namePublishedInYear,higherClassification,kingdom,phylum,class,order,family,genus,subgenus,specificEpithet,infraspecificEpithet,taxonRank,verbatimTaxonRank,scientificNameAuthorship,vernacularName,nomenclaturalCode,taxonomicStatus,nomenclaturalStatus,taxonRemarks"

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

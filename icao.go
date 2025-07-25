package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

const milCodeHeaderLen int = 2

var (
	errParseCSV  = errors.New("error parsing CSV")
	errHeaderLen = errors.New("unexpected header length")
)

type icaoAircraft struct {
	Class     string
	Engine    string
	ModelCode string
}

// getIcaoToAircraftMap returns an ICAO id to aircraft record mapping.
func getIcaoToAircraftMap() (map[string]icaoAircraft, error) {
	const icaoListPath string = "./data/ICAOList.csv"

	// Parse the CSV file
	icaoAircraftMap, err := parseIcaoCsvToMap(icaoListPath)
	if err != nil {
		return nil, fmt.Errorf("getIcaoToAircraftMap: %w", errParseCSV)
	}

	return icaoAircraftMap, nil
}

// parseIcaoCsvToMap reads a CSV file and parses it into a map ICAO -> aircraft spec.
func parseIcaoCsvToMap(filePath string) (map[string]icaoAircraft, error) {
	// Open the CSV file
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		return nil, fmt.Errorf("parseIcaoToCsvMap: failed to open file: %w", fileErr)
	}
	defer file.Close()
	// defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		log.Printf("failed to close file: %v", err)
	//	}
	// }(file) // Ensure the file is closed when the function exits

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	headers, headerErr := reader.Read()
	if headerErr != nil {
		return nil, fmt.Errorf("parseIcaoCsvToMap: failed to read header: %w", headerErr)
	}

	var icaoAircraftHeaders = [...]string{
		"aircraft TypeDesignator",
		"Class",
		"Number+Engine Type",
		"\"MANUFACTURER, Model\"",
	}
	if len(headers) != len(icaoAircraftHeaders) {
		return nil, fmt.Errorf("parseIcaoToCsvMap: %w", errHeaderLen)
	}

	records := make(map[string]icaoAircraft)

	// Loop through the remaining records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}

		if err != nil {
			return nil, fmt.Errorf("parseIcaoCsvToMap: failed to read record: %w", err)
		}

		key := record[0]
		class := record[1]
		engine := record[2]
		manufacturer := record[3]
		records[key] = icaoAircraft{class, engine, manufacturer}
	}

	return records, nil
}

// getMilCodeToOperatorMap returns a militar code to operator mapping.
func getMilCodeToOperatorMap() (map[string]string, error) {
	const milCodeFilePath string = "./data/MilICAOOperatorLookUp.csv"

	// Parse the CSV file
	icaoAircraftMap, err := parseMilCodeToMap(milCodeFilePath)
	if err != nil {
		return nil, fmt.Errorf("milCodeFilePath: %w", err)
	}

	return icaoAircraftMap, nil
}

// parseMilCodeToMap reads a CSV file and parses it into a map code -> military operator.
func parseMilCodeToMap(filePath string) (map[string]string, error) {
	// Open the CSV file
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		return nil, fmt.Errorf("parseMilCodeToMap: failed to open file: %w", fileErr)
	}
	defer file.Close()
	// defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		log.Printf("failed to close file: %v", err)
	//	}
	// }(file) // Ensure the file is closed when the function exits

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	headers, headerErr := reader.Read()
	if headerErr != nil {
		return nil, fmt.Errorf("parseMilCodeToMap: failed to read headers: %w", headerErr)
	}

	if len(headers) != milCodeHeaderLen {
		return nil, fmt.Errorf("parseMilCodeToMap: %w", errHeaderLen)
	}

	records := make(map[string]string)

	// Loop through the remaining records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}

		if err != nil {
			return nil, fmt.Errorf("parseMilCodeToMap: failed to read record: %w", err)
		}

		key := record[1]

		if len(key) == 0 {
			continue
		}

		militaryOperator := record[0]
		records[key] = militaryOperator
	}

	return records, nil
}

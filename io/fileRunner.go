////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package io

import (
	"encoding/csv"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"os"
	"strconv"
)

type AddressUpdate struct {
	Address string
	Value   uint64
}

// Creates an Addresses object from a CSV file
func StartIo(path string) (map[string]uint64, chan AddressUpdate) {
	// Open the file
	csvFile, err := os.Open(path)
	if err != nil {
		jww.FATAL.Panicf("Couldn't open the csv file: %+v", err)
	}
	defer csvFile.Close()

	// Parse the file
	r := csv.NewReader(csvFile)

	addressMap := make(map[string]uint64)

	// Iterate through the records
	i := 1
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			jww.FATAL.Panicf("Unable to read line: %+v", err)
		}
		if len(record) != 2 {
			jww.FATAL.Panicf("Line is formatted incorrectly: %v", record)
		}
		jww.DEBUG.Printf("Reading record %d: [Address: %s Amount: %s]", i, record[0], record[1])
		addressMap[record[0]], err = strconv.ParseUint(record[1], 10, 64)
		if err != nil {
			jww.FATAL.Panicf("Unable to process value: %+v", err)
		}
		i++
	}
	jww.INFO.Printf("Processed %d records!", len(addressMap))

	internalAddressMap := make(map[string]uint64)
	for key, value := range addressMap {
		internalAddressMap[key] = value
	}

	ch := make(chan AddressUpdate, 1000)
	go writeAddresses(path, internalAddressMap, ch)
	return addressMap, ch
}

// Writes Addresses object to CSV in an infinite loop
func writeAddresses(path string, addressMap map[string]uint64, ch chan AddressUpdate) {
	for {
		select {
		case newRecord := <-ch:
			jww.DEBUG.Printf("Updating record: %+v", newRecord)
			addressMap[newRecord.Address] = newRecord.Value

			// Open the file
			csvFile, err := os.OpenFile(path, os.O_WRONLY, 0755)
			if err != nil {
				jww.FATAL.Panicf("Couldn't open the csv file: %+v", err)
			}
			writer := csv.NewWriter(csvFile)

			i := 1
			records := make([][]string, 0)
			for key, value := range addressMap {
				record := make([]string, 2)
				record[0] = key
				record[1] = strconv.FormatUint(value, 10)
				records = append(records, record)
				jww.DEBUG.Printf("Writing record %d: %+v", i, record)
				i++
			}

			err = writer.WriteAll(records)
			if err != nil {
				jww.FATAL.Panicf("Unable to write to csv: %+v", err)
			}

			err = csvFile.Close()
			if err != nil {
				jww.FATAL.Panicf("Unable to close csv: %+v", err)
			}
		}
	}
}

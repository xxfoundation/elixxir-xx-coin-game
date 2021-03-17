////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package cmd

import (
	"encoding/csv"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Addresses struct {
	addressMap map[string]uint64
	sync.Mutex
}

// Modify the value of an address in the Addresses object
func (a *Addresses) Write(address string, value uint64) error {
	a.Lock()
	defer a.Unlock()

	if val, ok := a.addressMap[address]; ok {
		if val != 0 {
			return errors.Errorf("Address %s has already been processed", address)
		}
		jww.INFO.Printf("Updating value of %s to %d", address, value)
		a.addressMap[address] = val
		return nil
	}
	return errors.Errorf("Invalid address: %s", address)
}

// Creates an Addresses object from a CSV file
func readAddresses(path string) *Addresses {
	// Open the file
	csvFile, err := os.Open(path)
	if err != nil {
		jww.FATAL.Panicf("Couldn't open the csv file: %+v", err)
	}
	defer csvFile.Close()

	// Parse the file
	r := csv.NewReader(csvFile)

	addresses := &Addresses{
		addressMap: make(map[string]uint64),
	}

	// Iterate through the records
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
		jww.DEBUG.Printf("Address: %s, Amount: %s", record[0], record[1])
		addresses.addressMap[record[0]], err = strconv.ParseUint(record[1], 10, 64)
		if err != nil {
			jww.FATAL.Panicf("Unable to process value: %+v", err)
		}
	}
	jww.INFO.Printf("Processed %d records!", len(addresses.addressMap))
	return addresses
}

// Writes Addresses object to CSV in an infinite loop
func (a *Addresses) writeAddresses(path string) {
	for {
		// Open the file
		csvFile, err := os.OpenFile(path, os.O_WRONLY, 0755)
		if err != nil {
			jww.FATAL.Panicf("Couldn't open the csv file: %+v", err)
		}
		writer := csv.NewWriter(csvFile)

		i := 0
		records := make([][]string, 0)
		a.Lock()
		for key, value := range a.addressMap {
			record := make([]string, 2)
			record[0] = key
			record[1] = strconv.FormatUint(value, 10)
			records = append(records, record)
			jww.DEBUG.Printf("Writing record %d: %+v", i, record)
			i++
		}
		a.Unlock()

		err = writer.WriteAll(records)
		if err != nil {
			jww.FATAL.Panicf("Unable to write to csv: %+v", err)
		}

		err = csvFile.Close()
		if err != nil {
			jww.FATAL.Panicf("Unable to close csv: %+v", err)
		}
		time.Sleep(10 * time.Second)
	}
}
